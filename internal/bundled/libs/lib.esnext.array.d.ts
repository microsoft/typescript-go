/*! *****************************************************************************
Copyright (c) Microsoft Corporation. All rights reserved.
Licensed under the Apache License, Version 2.0 (the "License"); you may not use
this file except in compliance with the License. You may obtain a copy of the
License at http://www.apache.org/licenses/LICENSE-2.0

THIS CODE IS PROVIDED ON AN *AS IS* BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, EITHER EXPRESS OR IMPLIED, INCLUDING WITHOUT LIMITATION ANY IMPLIED
WARRANTIES OR CONDITIONS OF TITLE, FITNESS FOR A PARTICULAR PURPOSE,
MERCHANTABLITY OR NON-INFRINGEMENT.

See the Apache Version 2.0 License for specific language governing permissions
and limitations under the License.
***************************************************************************** */


/// <reference no-default-lib="true"/>

interface ArrayConstructor {
    /**
     * Creates an array from an async iterator or iterable object.
     * @param iterableOrArrayLike An async iterator or array-like object to convert to an array.
     */
    fromAsync<T>(iterableOrArrayLike: AsyncIterable<T> | Iterable<T | PromiseLike<T>> | ArrayLike<T | PromiseLike<T>>): Promise<T[]>;

    /**
     * Creates an array from an async iterator or iterable object.
     *
     * @param iterableOrArrayLike An async iterator or array-like object to convert to an array.
     * @param mapfn A mapping function to call on every element of itarableOrArrayLike.
     *      Each return value is awaited before being added to result array.
     * @param thisArg Value of 'this' used when executing mapfn.
     */
    fromAsync<T, U>(iterableOrArrayLike: AsyncIterable<T> | Iterable<T> | ArrayLike<T>, mapFn: (value: Awaited<T>, index: number) => U, thisArg?: any): Promise<Awaited<U>[]>;
}

interface Uint8ArrayConstructor {
    /**
     * Creates a new `Uint8Array` from a base64-encoded string.
     * @param string The base64-encoded string.
     * @param options If provided, specifies the alphabet and handling of the last chunk.
     * @returns A new `Uint8Array` instance.
     * @throws {SyntaxError} If the input string contains characters outside the specified alphabet, or if the last
     * chunk is inconsistent with the `lastChunkHandling` option.
     */
    fromBase64(
        string: string,
        options?: {
            alphabet?: "base64" | "base64url";
            lastChunkHandling?: "loose" | "strict" | "stop-before-partial";
        },
    ): Uint8Array<ArrayBuffer>;

    /**
     * Creates a new `Uint8Array` from a base16-encoded string.
     * @returns A new `Uint8Array` instance.
     */
    fromHex(
        string: string,
    ): Uint8Array<ArrayBuffer>;
}

interface Uint8Array<TArrayBuffer extends ArrayBufferLike> {
    /**
     * Converts the `Uint8Array` to a base64-encoded string.
     * @param options If provided, sets the alphabet and padding behavior used.
     * @returns A base64-encoded string.
     */
    toBase64(options?: {
        alphabet?: "base64" | "base64url";
        omitPadding?: boolean;
    }): string;

    /**
     * Sets the `Uint8Array` from a base64-encoded string.
     * @param string The base64-encoded string.
     * @param options If provided, specifies the alphabet and handling of the last chunk.
     * @returns An object containing the number of bytes read and written.
     * @throws {SyntaxError} If the input string contains characters outside the specified alphabet, or if the last
     * chunk is inconsistent with the `lastChunkHandling` option.
     */
    setFromBase64(string: string, options?: {
        alphabet?: "base64" | "base64url";
        lastChunkHandling?: "loose" | "strict" | "stop-before-partial";
    }): {
        read: number;
        written: number;
    };

    /**
     * Converts the `Uint8Array` to a base16-encoded string.
     * @returns A base16-encoded string.
     */
    toHex(): string;

    /**
     * Sets the `Uint8Array` from a base16-encoded string.
     * @param string The base16-encoded string.
     * @returns An object containing the number of bytes read and written.
     */
    setFromHex(string: string): {
        read: number;
        written: number;
    };
}
