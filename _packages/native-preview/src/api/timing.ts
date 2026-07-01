/**
 * Client-side collection of per-request timing and transfer measurements.
 *
 * When enabled, each request records its round-trip latency, the number of
 * payload bytes sent and received, and (when the transport reports it) the
 * server's own processing time. From those, an estimated transport overhead
 * (round-trip minus server processing time) is derived per request and across
 * running totals.
 */

/** Number of most-recent requests retained in the ring buffer. */
export const RECENT_REQUEST_CAPACITY = 5;

/** A single request's measured timing and transfer sample. */
export interface RequestTiming {
    /** The API method that was invoked. */
    method: string;
    /** Wall-clock round-trip time measured by the client, in milliseconds. */
    roundTripMs: number;
    /** Number of request payload bytes sent to the server. */
    bytesSent: number;
    /** Number of response payload bytes received from the server. */
    bytesReceived: number;
    /**
     * Server-side processing time in milliseconds, as reported by the transport.
     * It is `undefined` only in the rare case a response arrives without the expected
     * timing metadata (e.g. a malformed or empty response).
     */
    serverTimeMs: number | undefined;
    /**
     * Estimated transport overhead (round-trip minus server processing time) in
     * milliseconds. Only present when {@link serverTimeMs} is known. Clamped to
     * be non-negative.
     */
    transportOverheadMs: number | undefined;
    /** Wall-clock timestamp ({@link Date.now}) captured when the request completed. */
    timestamp: number;
}

/** Running totals accumulated across every measured request. */
export interface TimingAccumulators {
    /** Total number of requests measured. */
    requestCount: number;
    /** Sum of round-trip latencies, in milliseconds. */
    totalRoundTripMs: number;
    /** Total request payload bytes sent. */
    totalBytesSent: number;
    /** Total response payload bytes received. */
    totalBytesReceived: number;
    /** Sum of server processing time across requests that reported it, in milliseconds. */
    totalServerTimeMs: number;
    /**
     * Sum of estimated transport overhead across requests that reported server
     * time, in milliseconds.
     */
    totalTransportOverheadMs: number;
}

/** A point-in-time snapshot of collected timing information. */
export interface TimingInfo {
    /** Whether timing collection is enabled for this API instance. */
    enabled: boolean;
    /** Running totals across every measured request. */
    totals: TimingAccumulators;
    /**
     * The most recent requests, up to {@link RECENT_REQUEST_CAPACITY}, ordered
     * from oldest to newest.
     */
    recentRequests: RequestTiming[];
}

/** A raw measurement handed to {@link TimingCollector.record}. */
export interface TimingSample {
    method: string;
    roundTripMs: number;
    bytesSent: number;
    bytesReceived: number;
    /**
     * Server processing time in microseconds, as reported by the transport, or
     * `undefined` if the transport does not report it.
     */
    serverTimeMicros: number | undefined;
}

function emptyAccumulators(): TimingAccumulators {
    return {
        requestCount: 0,
        totalRoundTripMs: 0,
        totalBytesSent: 0,
        totalBytesReceived: 0,
        totalServerTimeMs: 0,
        totalTransportOverheadMs: 0,
    };
}

/** Returns a snapshot representing a disabled (never-collecting) timing state. */
export function disabledTimingInfo(): TimingInfo {
    return {
        enabled: false,
        totals: emptyAccumulators(),
        recentRequests: [],
    };
}

/**
 * Accumulates request timing samples into running totals and a fixed-size ring
 * buffer of the most recent requests.
 */
export class TimingCollector {
    private totals: TimingAccumulators = emptyAccumulators();
    // Ring buffer of the most recent requests. `ring` grows to at most
    // RECENT_REQUEST_CAPACITY; once full, `head` marks the oldest entry.
    private ring: RequestTiming[] = [];
    private head = 0;

    /** Records a single request's measurements. */
    record(sample: TimingSample): void {
        const serverTimeMs = sample.serverTimeMicros === undefined
            ? undefined
            : sample.serverTimeMicros / 1000;
        let transportOverheadMs: number | undefined;
        if (serverTimeMs !== undefined) {
            transportOverheadMs = Math.max(0, sample.roundTripMs - serverTimeMs);
            this.totals.totalServerTimeMs += serverTimeMs;
            this.totals.totalTransportOverheadMs += transportOverheadMs;
        }

        this.totals.requestCount++;
        this.totals.totalRoundTripMs += sample.roundTripMs;
        this.totals.totalBytesSent += sample.bytesSent;
        this.totals.totalBytesReceived += sample.bytesReceived;

        const entry: RequestTiming = {
            method: sample.method,
            roundTripMs: sample.roundTripMs,
            bytesSent: sample.bytesSent,
            bytesReceived: sample.bytesReceived,
            serverTimeMs,
            transportOverheadMs,
            timestamp: Date.now(),
        };

        if (this.ring.length < RECENT_REQUEST_CAPACITY) {
            this.ring.push(entry);
        }
        else {
            this.ring[this.head] = entry;
            this.head = (this.head + 1) % RECENT_REQUEST_CAPACITY;
        }
    }

    /** Returns a snapshot of the collected timing information. */
    getInfo(): TimingInfo {
        const recentRequests: RequestTiming[] = [];
        for (let i = 0; i < this.ring.length; i++) {
            recentRequests.push(this.ring[(this.head + i) % this.ring.length]);
        }
        return {
            enabled: true,
            totals: { ...this.totals },
            recentRequests,
        };
    }

    /** Clears all accumulated totals and recent-request history. */
    reset(): void {
        this.totals = emptyAccumulators();
        this.ring = [];
        this.head = 0;
    }
}
