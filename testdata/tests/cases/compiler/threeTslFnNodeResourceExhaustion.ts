// @strict: true
// @noEmit: true
// @skipLibCheck: true
// @noTypesAndSymbols: true
// @module: esnext
// @moduleResolution: bundler
// @target: esnext
// @lib: esnext,dom

// @filename: /src/repro.ts
import { Fn } from "three/tsl";
import type { Node } from "three/webgpu";

const hashValue = Fn(([cell]: [Node]) => cell);

void hashValue;

// @filename: /node_modules/three/package.json
{"name":"three","type":"module","exports":{"./tsl":"./build/three.tsl.js","./webgpu":"./build/three.webgpu.js"}}

// @filename: /node_modules/three/build/three.tsl.d.ts
export * from "../src/Three.TSL.js";


// @filename: /node_modules/three/build/three.webgpu.d.ts
export * from "../src/Three.WebGPU.js";


// @filename: /node_modules/three/src/Three.TSL.d.ts
import * as TSL from "./nodes/TSL.js";

export const BRDF_GGX: typeof TSL.BRDF_GGX;
export const BRDF_Lambert: typeof TSL.BRDF_Lambert;
export const BasicPointShadowFilter: typeof TSL.BasicPointShadowFilter;
export const BasicShadowFilter: typeof TSL.BasicShadowFilter;
export const Break: typeof TSL.Break;
export const Const: typeof TSL.Const;
export const Continue: typeof TSL.Continue;
export const DFGLUT: typeof TSL.DFGLUT;
export const D_GGX: typeof TSL.D_GGX;
export const Discard: typeof TSL.Discard;
export const EPSILON: typeof TSL.EPSILON;
export const F_Schlick: typeof TSL.F_Schlick;
export const Fn: typeof TSL.Fn;
export const INFINITY: typeof TSL.INFINITY;
export const If: typeof TSL.If;
export const Loop: typeof TSL.Loop;
export const NodeAccess: typeof TSL.NodeAccess;
export const NodeShaderStage: typeof TSL.NodeShaderStage;
export const NodeType: typeof TSL.NodeType;
export const NodeUpdateType: typeof TSL.NodeUpdateType;
export const PCFShadowFilter: typeof TSL.PCFShadowFilter;
export const PCFSoftShadowFilter: typeof TSL.PCFSoftShadowFilter;
export const PI: typeof TSL.PI;
export const PI2: typeof TSL.PI2;
export const TWO_PI: typeof TSL.TWO_PI;
export const HALF_PI: typeof TSL.HALF_PI;
export const PointShadowFilter: typeof TSL.PointShadowFilter;
export const Return: typeof TSL.Return;
export const Schlick_to_F0: typeof TSL.Schlick_to_F0;
export const ShaderNode: typeof TSL.ShaderNode;
export const Stack: typeof TSL.Stack;
export const Switch: typeof TSL.Switch;
export const TBNViewMatrix: typeof TSL.TBNViewMatrix;
export const VSMShadowFilter: typeof TSL.VSMShadowFilter;
export const V_GGX_SmithCorrelated: typeof TSL.V_GGX_SmithCorrelated;
export const Var: typeof TSL.Var;
export const VarIntent: typeof TSL.VarIntent;
export const abs: typeof TSL.abs;
export const acesFilmicToneMapping: typeof TSL.acesFilmicToneMapping;
export const acos: typeof TSL.acos;
export const acosh: typeof TSL.acosh;
export const add: typeof TSL.add;
export const addMethodChaining: typeof TSL.addMethodChaining;
export const addNodeElement: typeof TSL.addNodeElement;
export const agxToneMapping: typeof TSL.agxToneMapping;
export const all: typeof TSL.all;
export const alphaT: typeof TSL.alphaT;
export const ambientOcclusion: typeof TSL.ambientOcclusion;
export const and: typeof TSL.and;
export const anisotropy: typeof TSL.anisotropy;
export const anisotropyB: typeof TSL.anisotropyB;
export const anisotropyT: typeof TSL.anisotropyT;
export const any: typeof TSL.any;
export const append: typeof TSL.append;
export const array: typeof TSL.array;
export const asin: typeof TSL.asin;
export const asinh: typeof TSL.asinh;
export const assign: typeof TSL.assign;
export const atan: typeof TSL.atan;
export const atanh: typeof TSL.atanh;
export const atomicAdd: typeof TSL.atomicAdd;
export const atomicAnd: typeof TSL.atomicAnd;
export const atomicFunc: typeof TSL.atomicFunc;
export const atomicLoad: typeof TSL.atomicLoad;
export const atomicMax: typeof TSL.atomicMax;
export const atomicMin: typeof TSL.atomicMin;
export const atomicOr: typeof TSL.atomicOr;
export const atomicStore: typeof TSL.atomicStore;
export const atomicSub: typeof TSL.atomicSub;
export const atomicXor: typeof TSL.atomicXor;
export const attenuationColor: typeof TSL.attenuationColor;
export const attenuationDistance: typeof TSL.attenuationDistance;
export const attribute: typeof TSL.attribute;
export const attributeArray: typeof TSL.attributeArray;
export const backgroundBlurriness: typeof TSL.backgroundBlurriness;
export const backgroundIntensity: typeof TSL.backgroundIntensity;
export const backgroundRotation: typeof TSL.backgroundRotation;
export const batch: typeof TSL.batch;
export const bentNormalView: typeof TSL.bentNormalView;
export const billboarding: typeof TSL.billboarding;
export const bitAnd: typeof TSL.bitAnd;
export const bitNot: typeof TSL.bitNot;
export const bitOr: typeof TSL.bitOr;
export const bitXor: typeof TSL.bitXor;
export const bitangentGeometry: typeof TSL.bitangentGeometry;
export const bitangentLocal: typeof TSL.bitangentLocal;
export const bitangentView: typeof TSL.bitangentView;
export const bitangentWorld: typeof TSL.bitangentWorld;
export const bitcast: typeof TSL.bitcast;
export const blendBurn: typeof TSL.blendBurn;
export const blendColor: typeof TSL.blendColor;
export const blendDodge: typeof TSL.blendDodge;
export const blendOverlay: typeof TSL.blendOverlay;
export const blendScreen: typeof TSL.blendScreen;
export const blur: typeof TSL.blur;
export const bool: typeof TSL.bool;
export const buffer: typeof TSL.buffer;
export const bufferAttribute: typeof TSL.bufferAttribute;
export const bumpMap: typeof TSL.bumpMap;
export const builtin: typeof TSL.builtin;
export const builtinAOContext: typeof TSL.builtinAOContext;
export const builtinShadowContext: typeof TSL.builtinShadowContext;
export const bvec2: typeof TSL.bvec2;
export const bvec3: typeof TSL.bvec3;
export const bvec4: typeof TSL.bvec4;
export const bypass: typeof TSL.bypass;
export const cache: typeof TSL.cache;
export const call: typeof TSL.call;
export const cameraFar: typeof TSL.cameraFar;
export const cameraIndex: typeof TSL.cameraIndex;
export const cameraNear: typeof TSL.cameraNear;
export const cameraNormalMatrix: typeof TSL.cameraNormalMatrix;
export const cameraPosition: typeof TSL.cameraPosition;
export const cameraProjectionMatrix: typeof TSL.cameraProjectionMatrix;
export const cameraProjectionMatrixInverse: typeof TSL.cameraProjectionMatrixInverse;
export const cameraViewMatrix: typeof TSL.cameraViewMatrix;
export const cameraViewport: typeof TSL.cameraViewport;
export const cameraWorldMatrix: typeof TSL.cameraWorldMatrix;
export const cbrt: typeof TSL.cbrt;
export const cdl: typeof TSL.cdl;
export const ceil: typeof TSL.ceil;
export const checker: typeof TSL.checker;
export const cineonToneMapping: typeof TSL.cineonToneMapping;
export const clamp: typeof TSL.clamp;
export const clearcoat: typeof TSL.clearcoat;
export const clearcoatNormalView: typeof TSL.clearcoatNormalView;
export const clearcoatRoughness: typeof TSL.clearcoatRoughness;
export const clipSpace: typeof TSL.clipSpace;
export const code: typeof TSL.code;
export const color: typeof TSL.color;
export const colorSpaceToWorking: typeof TSL.colorSpaceToWorking;
export const colorToDirection: typeof TSL.colorToDirection;
export const compute: typeof TSL.compute;
export const computeKernel: typeof TSL.computeKernel;
export const computeSkinning: typeof TSL.computeSkinning;
export const context: typeof TSL.context;
export const convert: typeof TSL.convert;
export const convertColorSpace: typeof TSL.convertColorSpace;
export const convertToTexture: typeof TSL.convertToTexture;
export const countLeadingZeros: typeof TSL.countLeadingZeros;
export const countOneBits: typeof TSL.countOneBits;
export const countTrailingZeros: typeof TSL.countTrailingZeros;
export const cos: typeof TSL.cos;
export const cosh: typeof TSL.cosh;
export const cross: typeof TSL.cross;
export const cubeTexture: typeof TSL.cubeTexture;
export const cubeTextureBase: typeof TSL.cubeTextureBase;
export const dFdx: typeof TSL.dFdx;
export const dFdy: typeof TSL.dFdy;
export const dashSize: typeof TSL.dashSize;
export const debug: typeof TSL.debug;
export const decrement: typeof TSL.decrement;
export const decrementBefore: typeof TSL.decrementBefore;
export const defaultBuildStages: typeof TSL.defaultBuildStages;
export const defaultShaderStages: typeof TSL.defaultShaderStages;
export const defined: typeof TSL.defined;
export const degrees: typeof TSL.degrees;
export const deltaTime: typeof TSL.deltaTime;
export const densityFog: typeof TSL.densityFog;
export const densityFogFactor: typeof TSL.densityFogFactor;
export const depth: typeof TSL.depth;
export const depthPass: typeof TSL.depthPass;
export const determinant: typeof TSL.determinant;
export const difference: typeof TSL.difference;
export const diffuseColor: typeof TSL.diffuseColor;
export const directPointLight: typeof TSL.directPointLight;
export const directionToColor: typeof TSL.directionToColor;
export const directionToFaceDirection: typeof TSL.directionToFaceDirection;
export const dispersion: typeof TSL.dispersion;
export const distance: typeof TSL.distance;
export const div: typeof TSL.div;
export const dot: typeof TSL.dot;
export const drawIndex: typeof TSL.drawIndex;
export const dynamicBufferAttribute: typeof TSL.dynamicBufferAttribute;
export const element: typeof TSL.element;
export const emissive: typeof TSL.emissive;
export const equal: typeof TSL.equal;
export const equirectDirection: typeof TSL.equirectDirection;
export const equirectUV: typeof TSL.equirectUV;
export const exp: typeof TSL.exp;
export const exp2: typeof TSL.exp2;
export const exponentialHeightFogFactor: typeof TSL.exponentialHeightFogFactor;
export const expression: typeof TSL.expression;
export const faceDirection: typeof TSL.faceDirection;
export const faceForward: typeof TSL.faceForward;
export const faceforward: typeof TSL.faceforward;
export const float: typeof TSL.float;
export const floatBitsToInt: typeof TSL.floatBitsToInt;
export const floatBitsToUint: typeof TSL.floatBitsToUint;
export const floor: typeof TSL.floor;
export const fog: typeof TSL.fog;
export const fract: typeof TSL.fract;
export const frameGroup: typeof TSL.frameGroup;
export const frameId: typeof TSL.frameId;
export const frontFacing: typeof TSL.frontFacing;
export const fwidth: typeof TSL.fwidth;
export const gain: typeof TSL.gain;
export const gapSize: typeof TSL.gapSize;
export const getConstNodeType: typeof TSL.getConstNodeType;
export const getCurrentStack: typeof TSL.getCurrentStack;
export const getDirection: typeof TSL.getDirection;
export const getDistanceAttenuation: typeof TSL.getDistanceAttenuation;
export const getGeometryRoughness: typeof TSL.getGeometryRoughness;
export const getNormalFromDepth: typeof TSL.getNormalFromDepth;
export const interleavedGradientNoise: typeof TSL.interleavedGradientNoise;
export const vogelDiskSample: typeof TSL.vogelDiskSample;
export const getParallaxCorrectNormal: typeof TSL.getParallaxCorrectNormal;
export const getRoughness: typeof TSL.getRoughness;
export const getScreenPosition: typeof TSL.getScreenPosition;
export const getShIrradianceAt: typeof TSL.getShIrradianceAt;
export const getShadowMaterial: typeof TSL.getShadowMaterial;
export const getShadowRenderObjectFunction: typeof TSL.getShadowRenderObjectFunction;
export const getTextureIndex: typeof TSL.getTextureIndex;
export const getViewPosition: typeof TSL.getViewPosition;
export const globalId: typeof TSL.globalId;
export const glsl: typeof TSL.glsl;
export const glslFn: typeof TSL.glslFn;
export const grayscale: typeof TSL.grayscale;
export const greaterThan: typeof TSL.greaterThan;
export const greaterThanEqual: typeof TSL.greaterThanEqual;
export const hash: typeof TSL.hash;
export const highpModelNormalViewMatrix: typeof TSL.highpModelNormalViewMatrix;
export const highpModelViewMatrix: typeof TSL.highpModelViewMatrix;
export const hue: typeof TSL.hue;
export const increment: typeof TSL.increment;
export const incrementBefore: typeof TSL.incrementBefore;
export const instance: typeof TSL.instance;
export const instanceIndex: typeof TSL.instanceIndex;
export const instancedArray: typeof TSL.instancedArray;
export const instancedBufferAttribute: typeof TSL.instancedBufferAttribute;
export const instancedDynamicBufferAttribute: typeof TSL.instancedDynamicBufferAttribute;
export const instancedMesh: typeof TSL.instancedMesh;
export const int: typeof TSL.int;
export const intBitsToFloat: typeof TSL.intBitsToFloat;
export const inverse: typeof TSL.inverse;
export const inverseSqrt: typeof TSL.inverseSqrt;
export const inversesqrt: typeof TSL.inversesqrt;
export const invocationLocalIndex: typeof TSL.invocationLocalIndex;
export const invocationSubgroupIndex: typeof TSL.invocationSubgroupIndex;
export const ior: typeof TSL.ior;
export const iridescence: typeof TSL.iridescence;
export const iridescenceIOR: typeof TSL.iridescenceIOR;
export const iridescenceThickness: typeof TSL.iridescenceThickness;
export const ivec2: typeof TSL.ivec2;
export const ivec3: typeof TSL.ivec3;
export const ivec4: typeof TSL.ivec4;
export const js: typeof TSL.js;
export const label: typeof TSL.label;
export const length: typeof TSL.length;
export const lengthSq: typeof TSL.lengthSq;
export const lessThan: typeof TSL.lessThan;
export const lessThanEqual: typeof TSL.lessThanEqual;
export const lightPosition: typeof TSL.lightPosition;
export const lightProjectionUV: typeof TSL.lightProjectionUV;
export const lightShadowMatrix: typeof TSL.lightShadowMatrix;
export const lightTargetDirection: typeof TSL.lightTargetDirection;
export const lightTargetPosition: typeof TSL.lightTargetPosition;
export const lightViewPosition: typeof TSL.lightViewPosition;
export const lightingContext: typeof TSL.lightingContext;
export const lights: typeof TSL.lights;
export const linearDepth: typeof TSL.linearDepth;
export const linearToneMapping: typeof TSL.linearToneMapping;
export const localId: typeof TSL.localId;
export const log: typeof TSL.log;
export const log2: typeof TSL.log2;
export const logarithmicDepthToViewZ: typeof TSL.logarithmicDepthToViewZ;
export const luminance: typeof TSL.luminance;
export const mat2: typeof TSL.mat2;
export const mat3: typeof TSL.mat3;
export const mat4: typeof TSL.mat4;
export const matcapUV: typeof TSL.matcapUV;
export const materialAO: typeof TSL.materialAO;
export const materialAlphaTest: typeof TSL.materialAlphaTest;
export const materialAnisotropy: typeof TSL.materialAnisotropy;
export const materialAnisotropyVector: typeof TSL.materialAnisotropyVector;
export const materialAttenuationColor: typeof TSL.materialAttenuationColor;
export const materialAttenuationDistance: typeof TSL.materialAttenuationDistance;
export const materialClearcoat: typeof TSL.materialClearcoat;
export const materialClearcoatNormal: typeof TSL.materialClearcoatNormal;
export const materialClearcoatRoughness: typeof TSL.materialClearcoatRoughness;
export const materialColor: typeof TSL.materialColor;
export const materialDispersion: typeof TSL.materialDispersion;
export const materialEmissive: typeof TSL.materialEmissive;
export const materialEnvIntensity: typeof TSL.materialEnvIntensity;
export const materialEnvRotation: typeof TSL.materialEnvRotation;
export const materialIOR: typeof TSL.materialIOR;
export const materialIridescence: typeof TSL.materialIridescence;
export const materialIridescenceIOR: typeof TSL.materialIridescenceIOR;
export const materialIridescenceThickness: typeof TSL.materialIridescenceThickness;
export const materialLightMap: typeof TSL.materialLightMap;
export const materialLineDashOffset: typeof TSL.materialLineDashOffset;
export const materialLineDashSize: typeof TSL.materialLineDashSize;
export const materialLineGapSize: typeof TSL.materialLineGapSize;
export const materialLineScale: typeof TSL.materialLineScale;
export const materialLineWidth: typeof TSL.materialLineWidth;
export const materialMetalness: typeof TSL.materialMetalness;
export const materialNormal: typeof TSL.materialNormal;
export const materialOpacity: typeof TSL.materialOpacity;
export const materialPointSize: typeof TSL.materialPointSize;
export const materialReference: typeof TSL.materialReference;
export const materialReflectivity: typeof TSL.materialReflectivity;
export const materialRefractionRatio: typeof TSL.materialRefractionRatio;
export const materialRotation: typeof TSL.materialRotation;
export const materialRoughness: typeof TSL.materialRoughness;
export const materialSheen: typeof TSL.materialSheen;
export const materialSheenRoughness: typeof TSL.materialSheenRoughness;
export const materialShininess: typeof TSL.materialShininess;
export const materialSpecular: typeof TSL.materialSpecular;
export const materialSpecularColor: typeof TSL.materialSpecularColor;
export const materialSpecularIntensity: typeof TSL.materialSpecularIntensity;
export const materialSpecularStrength: typeof TSL.materialSpecularStrength;
export const materialThickness: typeof TSL.materialThickness;
export const materialTransmission: typeof TSL.materialTransmission;
export const max: typeof TSL.max;
export const maxMipLevel: typeof TSL.maxMipLevel;
export const mediumpModelViewMatrix: typeof TSL.mediumpModelViewMatrix;
export const metalness: typeof TSL.metalness;
export const min: typeof TSL.min;
export const mix: typeof TSL.mix;
export const mixElement: typeof TSL.mixElement;
export const mod: typeof TSL.mod;
export const modelDirection: typeof TSL.modelDirection;
export const modelNormalMatrix: typeof TSL.modelNormalMatrix;
export const modelPosition: typeof TSL.modelPosition;
export const modelRadius: typeof TSL.modelRadius;
export const modelScale: typeof TSL.modelScale;
export const modelViewMatrix: typeof TSL.modelViewMatrix;
export const modelViewPosition: typeof TSL.modelViewPosition;
export const modelViewProjection: typeof TSL.modelViewProjection;
export const modelWorldMatrix: typeof TSL.modelWorldMatrix;
export const modelWorldMatrixInverse: typeof TSL.modelWorldMatrixInverse;
export const morphReference: typeof TSL.morphReference;
export const mrt: typeof TSL.mrt;
export const mul: typeof TSL.mul;
export const mx_aastep: typeof TSL.mx_aastep;
export const mx_add: typeof TSL.mx_add;
export const mx_atan2: typeof TSL.mx_atan2;
export const mx_cell_noise_float: typeof TSL.mx_cell_noise_float;
export const mx_contrast: typeof TSL.mx_contrast;
export const mx_divide: typeof TSL.mx_divide;
export const mx_fractal_noise_float: typeof TSL.mx_fractal_noise_float;
export const mx_fractal_noise_vec2: typeof TSL.mx_fractal_noise_vec2;
export const mx_fractal_noise_vec3: typeof TSL.mx_fractal_noise_vec3;
export const mx_fractal_noise_vec4: typeof TSL.mx_fractal_noise_vec4;
export const mx_frame: typeof TSL.mx_frame;
export const mx_heighttonormal: typeof TSL.mx_heighttonormal;
export const mx_hsvtorgb: typeof TSL.mx_hsvtorgb;
export const mx_ifequal: typeof TSL.mx_ifequal;
export const mx_ifgreater: typeof TSL.mx_ifgreater;
export const mx_ifgreatereq: typeof TSL.mx_ifgreatereq;
export const mx_invert: typeof TSL.mx_invert;
export const mx_modulo: typeof TSL.mx_modulo;
export const mx_multiply: typeof TSL.mx_multiply;
export const mx_noise_float: typeof TSL.mx_noise_float;
export const mx_noise_vec3: typeof TSL.mx_noise_vec3;
export const mx_noise_vec4: typeof TSL.mx_noise_vec4;
export const mx_place2d: typeof TSL.mx_place2d;
export const mx_power: typeof TSL.mx_power;
export const mx_ramp4: typeof TSL.mx_ramp4;
export const mx_ramplr: typeof TSL.mx_ramplr;
export const mx_ramptb: typeof TSL.mx_ramptb;
export const mx_rgbtohsv: typeof TSL.mx_rgbtohsv;
export const mx_rotate2d: typeof TSL.mx_rotate2d;
export const mx_rotate3d: typeof TSL.mx_rotate3d;
export const mx_safepower: typeof TSL.mx_safepower;
export const mx_separate: typeof TSL.mx_separate;
export const mx_splitlr: typeof TSL.mx_splitlr;
export const mx_splittb: typeof TSL.mx_splittb;
export const mx_srgb_texture_to_lin_rec709: typeof TSL.mx_srgb_texture_to_lin_rec709;
export const mx_subtract: typeof TSL.mx_subtract;
export const mx_timer: typeof TSL.mx_timer;
export const mx_transform_uv: typeof TSL.mx_transform_uv;
export const mx_unifiednoise2d: typeof TSL.mx_unifiednoise2d;
export const mx_unifiednoise3d: typeof TSL.mx_unifiednoise3d;
export const mx_worley_noise_float: typeof TSL.mx_worley_noise_float;
export const mx_worley_noise_vec2: typeof TSL.mx_worley_noise_vec2;
export const mx_worley_noise_vec3: typeof TSL.mx_worley_noise_vec3;
export const negate: typeof TSL.negate;
export const negateOnBackSide: typeof TSL.negateOnBackSide;
export const neutralToneMapping: typeof TSL.neutralToneMapping;
export const nodeArray: typeof TSL.nodeArray;
export const nodeImmutable: typeof TSL.nodeImmutable;
export const nodeObject: typeof TSL.nodeObject;
export const nodeObjectIntent: typeof TSL.nodeObjectIntent;
export const nodeObjects: typeof TSL.nodeObjects;
export const nodeProxy: typeof TSL.nodeProxy;
export const nodeProxyIntent: typeof TSL.nodeProxyIntent;
export const normalFlat: typeof TSL.normalFlat;
export const normalGeometry: typeof TSL.normalGeometry;
export const normalLocal: typeof TSL.normalLocal;
export const normalMap: typeof TSL.normalMap;
export const normalView: typeof TSL.normalView;
export const normalViewGeometry: typeof TSL.normalViewGeometry;
export const normalWorld: typeof TSL.normalWorld;
export const normalWorldGeometry: typeof TSL.normalWorldGeometry;
export const normalize: typeof TSL.normalize;
export const not: typeof TSL.not;
export const notEqual: typeof TSL.notEqual;
export const numWorkgroups: typeof TSL.numWorkgroups;
export const objectDirection: typeof TSL.objectDirection;
export const objectGroup: typeof TSL.objectGroup;
export const objectPosition: typeof TSL.objectPosition;
export const objectRadius: typeof TSL.objectRadius;
export const objectScale: typeof TSL.objectScale;
export const objectViewPosition: typeof TSL.objectViewPosition;
export const objectWorldMatrix: typeof TSL.objectWorldMatrix;
export const OnBeforeObjectUpdate: typeof TSL.OnBeforeObjectUpdate;
export const OnBeforeMaterialUpdate: typeof TSL.OnBeforeMaterialUpdate;
export const OnObjectUpdate: typeof TSL.OnObjectUpdate;
export const OnMaterialUpdate: typeof TSL.OnMaterialUpdate;
export const oneMinus: typeof TSL.oneMinus;
export const or: typeof TSL.or;
export const orthographicDepthToViewZ: typeof TSL.orthographicDepthToViewZ;
export const oscSawtooth: typeof TSL.oscSawtooth;
export const oscSine: typeof TSL.oscSine;
export const oscSquare: typeof TSL.oscSquare;
export const oscTriangle: typeof TSL.oscTriangle;
export const output: typeof TSL.output;
export const outputStruct: typeof TSL.outputStruct;
export const overloadingFn: typeof TSL.overloadingFn;
export const overrideNode: typeof TSL.overrideNode;
export const overrideNodes: typeof TSL.overrideNodes;
export const packHalf2x16: typeof TSL.packHalf2x16;
export const packSnorm2x16: typeof TSL.packSnorm2x16;
export const packUnorm2x16: typeof TSL.packUnorm2x16;
export const packNormalToRGB: typeof TSL.packNormalToRGB;
export const parabola: typeof TSL.parabola;
export const parallaxDirection: typeof TSL.parallaxDirection;
export const parallaxUV: typeof TSL.parallaxUV;
export const parameter: typeof TSL.parameter;
export const pass: typeof TSL.pass;
export const passTexture: typeof TSL.passTexture;
export const pcurve: typeof TSL.pcurve;
export const perspectiveDepthToViewZ: typeof TSL.perspectiveDepthToViewZ;
export const pmremTexture: typeof TSL.pmremTexture;
export const pointShadow: typeof TSL.pointShadow;
export const pointUV: typeof TSL.pointUV;
export const pointWidth: typeof TSL.pointWidth;
export const positionGeometry: typeof TSL.positionGeometry;
export const positionLocal: typeof TSL.positionLocal;
export const positionPrevious: typeof TSL.positionPrevious;
export const positionView: typeof TSL.positionView;
export const positionViewDirection: typeof TSL.positionViewDirection;
export const positionWorld: typeof TSL.positionWorld;
export const positionWorldDirection: typeof TSL.positionWorldDirection;
export const posterize: typeof TSL.posterize;
export const pow: typeof TSL.pow;
export const pow2: typeof TSL.pow2;
export const pow3: typeof TSL.pow3;
export const pow4: typeof TSL.pow4;
export const premultiplyAlpha: typeof TSL.premultiplyAlpha;
export const property: typeof TSL.property;
export const radians: typeof TSL.radians;
export const rand: typeof TSL.rand;
export const range: typeof TSL.range;
export const rangeFog: typeof TSL.rangeFog;
export const rangeFogFactor: typeof TSL.rangeFogFactor;
export const reciprocal: typeof TSL.reciprocal;
export const reference: typeof TSL.reference;
export const referenceBuffer: typeof TSL.referenceBuffer;
export const reflect: typeof TSL.reflect;
export const reflectVector: typeof TSL.reflectVector;
export const reflectView: typeof TSL.reflectView;
export const reflector: typeof TSL.reflector;
export const refract: typeof TSL.refract;
export const refractVector: typeof TSL.refractVector;
export const refractView: typeof TSL.refractView;
export const reinhardToneMapping: typeof TSL.reinhardToneMapping;
export const remap: typeof TSL.remap;
export const remapClamp: typeof TSL.remapClamp;
export const renderGroup: typeof TSL.renderGroup;
export const renderOutput: typeof TSL.renderOutput;
export const rendererReference: typeof TSL.rendererReference;
export const replaceDefaultUV: typeof TSL.replaceDefaultUV;
export const rotate: typeof TSL.rotate;
export const rotateUV: typeof TSL.rotateUV;
export const roughness: typeof TSL.roughness;
export const round: typeof TSL.round;
export const rtt: typeof TSL.rtt;
export const sRGBTransferEOTF: typeof TSL.sRGBTransferEOTF;
export const sRGBTransferOETF: typeof TSL.sRGBTransferOETF;
export const sample: typeof TSL.sample;
export const sampler: typeof TSL.sampler;
export const samplerComparison: typeof TSL.samplerComparison;
export const saturate: typeof TSL.saturate;
export const saturation: typeof TSL.saturation;
export const screenCoordinate: typeof TSL.screenCoordinate;
export const screenDPR: typeof TSL.screenDPR;
export const screenSize: typeof TSL.screenSize;
export const screenUV: typeof TSL.screenUV;
export const select: typeof TSL.select;
export const setCurrentStack: typeof TSL.setCurrentStack;
export const setName: typeof TSL.setName;
export const shaderStages: typeof TSL.shaderStages;
export const shadow: typeof TSL.shadow;
export const shadowPositionWorld: typeof TSL.shadowPositionWorld;
export const shapeCircle: typeof TSL.shapeCircle;
export const sharedUniformGroup: typeof TSL.sharedUniformGroup;
export const sheen: typeof TSL.sheen;
export const sheenRoughness: typeof TSL.sheenRoughness;
export const shiftLeft: typeof TSL.shiftLeft;
export const shiftRight: typeof TSL.shiftRight;
export const shininess: typeof TSL.shininess;
export const sign: typeof TSL.sign;
export const sin: typeof TSL.sin;
export const sinh: typeof TSL.sinh;
export const sinc: typeof TSL.sinc;
export const skinning: typeof TSL.skinning;
export const smoothstep: typeof TSL.smoothstep;
export const smoothstepElement: typeof TSL.smoothstepElement;
export const specularColor: typeof TSL.specularColor;
export const specularF90: typeof TSL.specularF90;
export const spherizeUV: typeof TSL.spherizeUV;
export const split: typeof TSL.split;
export const spritesheetUV: typeof TSL.spritesheetUV;
export const sqrt: typeof TSL.sqrt;
export const stack: typeof TSL.stack;
export const step: typeof TSL.step;
export const stepElement: typeof TSL.stepElement;
export const storage: typeof TSL.storage;
export const storageBarrier: typeof TSL.storageBarrier;
export const storageTexture: typeof TSL.storageTexture;
export const storageTexture3D: typeof TSL.storageTexture3D;
export const struct: typeof TSL.struct;
export const sub: typeof TSL.sub;
export const subgroupAdd: typeof TSL.subgroupAdd;
export const subgroupAll: typeof TSL.subgroupAll;
export const subgroupAnd: typeof TSL.subgroupAnd;
export const subgroupAny: typeof TSL.subgroupAny;
export const subgroupBallot: typeof TSL.subgroupBallot;
export const subgroupBroadcast: typeof TSL.subgroupBroadcast;
export const subgroupBroadcastFirst: typeof TSL.subgroupBroadcastFirst;
export const subBuild: typeof TSL.subBuild;
export const subgroupElect: typeof TSL.subgroupElect;
export const subgroupExclusiveAdd: typeof TSL.subgroupExclusiveAdd;
export const subgroupExclusiveMul: typeof TSL.subgroupExclusiveMul;
export const subgroupInclusiveAdd: typeof TSL.subgroupInclusiveAdd;
export const subgroupInclusiveMul: typeof TSL.subgroupInclusiveMul;
export const subgroupIndex: typeof TSL.subgroupIndex;
export const subgroupMax: typeof TSL.subgroupMax;
export const subgroupMin: typeof TSL.subgroupMin;
export const subgroupMul: typeof TSL.subgroupMul;
export const subgroupOr: typeof TSL.subgroupOr;
export const subgroupShuffle: typeof TSL.subgroupShuffle;
export const subgroupShuffleDown: typeof TSL.subgroupShuffleDown;
export const subgroupShuffleUp: typeof TSL.subgroupShuffleUp;
export const subgroupShuffleXor: typeof TSL.subgroupShuffleXor;
export const subgroupSize: typeof TSL.subgroupSize;
export const subgroupXor: typeof TSL.subgroupXor;
export const tan: typeof TSL.tan;
export const tanh: typeof TSL.tanh;
export const tangentGeometry: typeof TSL.tangentGeometry;
export const tangentLocal: typeof TSL.tangentLocal;
export const tangentView: typeof TSL.tangentView;
export const tangentWorld: typeof TSL.tangentWorld;
export const texture: typeof TSL.texture;
export const texture3D: typeof TSL.texture3D;
export const textureBarrier: typeof TSL.textureBarrier;
export const textureBicubic: typeof TSL.textureBicubic;
export const textureBicubicLevel: typeof TSL.textureBicubicLevel;
export const textureCubeUV: typeof TSL.textureCubeUV;
export const textureLoad: typeof TSL.textureLoad;
export const textureSize: typeof TSL.textureSize;
export const textureLevel: typeof TSL.textureLevel;
export const textureStore: typeof TSL.textureStore;
export const thickness: typeof TSL.thickness;
export const time: typeof TSL.time;
export const toneMapping: typeof TSL.toneMapping;
export const toneMappingExposure: typeof TSL.toneMappingExposure;
export const toonOutlinePass: typeof TSL.toonOutlinePass;
export const transformDirection: typeof TSL.transformDirection;
export const transformNormal: typeof TSL.transformNormal;
export const transformNormalByInverseViewMatrix: typeof TSL.transformNormalByInverseViewMatrix;
export const transformNormalByViewMatrix: typeof TSL.transformNormalByViewMatrix;
export const transformNormalToView: typeof TSL.transformNormalToView;
export const transformedClearcoatNormalView: typeof TSL.transformedClearcoatNormalView;
export const transformedNormalView: typeof TSL.transformedNormalView;
export const transformedNormalWorld: typeof TSL.transformedNormalWorld;
export const transmission: typeof TSL.transmission;
export const transpose: typeof TSL.transpose;
export const triNoise3D: typeof TSL.triNoise3D;
export const triplanarTexture: typeof TSL.triplanarTexture;
export const triplanarTextures: typeof TSL.triplanarTextures;
export const trunc: typeof TSL.trunc;
export const uint: typeof TSL.uint;
export const uintBitsToFloat: typeof TSL.uintBitsToFloat;
export const uniform: typeof TSL.uniform;
export const uniformArray: typeof TSL.uniformArray;
export const uniformCubeTexture: typeof TSL.uniformCubeTexture;
export const uniformGroup: typeof TSL.uniformGroup;
export const uniformFlow: typeof TSL.uniformFlow;
export const uniformTexture: typeof TSL.uniformTexture;
export const unpackHalf2x16: typeof TSL.unpackHalf2x16;
export const unpackSnorm2x16: typeof TSL.unpackSnorm2x16;
export const unpackUnorm2x16: typeof TSL.unpackUnorm2x16;
export const unpackRGBToNormal: typeof TSL.unpackRGBToNormal;
export const unpremultiplyAlpha: typeof TSL.unpremultiplyAlpha;
export const userData: typeof TSL.userData;
export const uv: typeof TSL.uv;
export const uvec2: typeof TSL.uvec2;
export const uvec3: typeof TSL.uvec3;
export const uvec4: typeof TSL.uvec4;
export const varying: typeof TSL.varying;
export const varyingProperty: typeof TSL.varyingProperty;
export const vec2: typeof TSL.vec2;
export const vec3: typeof TSL.vec3;
export const vec4: typeof TSL.vec4;
export const vectorComponents: typeof TSL.vectorComponents;
export const velocity: typeof TSL.velocity;
export const vertexColor: typeof TSL.vertexColor;
export const vertexIndex: typeof TSL.vertexIndex;
export const vertexStage: typeof TSL.vertexStage;
export const vibrance: typeof TSL.vibrance;
export const viewZToLogarithmicDepth: typeof TSL.viewZToLogarithmicDepth;
export const viewZToOrthographicDepth: typeof TSL.viewZToOrthographicDepth;
export const viewZToPerspectiveDepth: typeof TSL.viewZToPerspectiveDepth;
export const viewZToReversedOrthographicDepth: typeof TSL.viewZToReversedOrthographicDepth;
export const viewZToReversedPerspectiveDepth: typeof TSL.viewZToReversedPerspectiveDepth;
export const viewport: typeof TSL.viewport;
export const viewportCoordinate: typeof TSL.viewportCoordinate;
export const viewportDepthTexture: typeof TSL.viewportDepthTexture;
export const viewportLinearDepth: typeof TSL.viewportLinearDepth;
export const viewportMipTexture: typeof TSL.viewportMipTexture;
export const viewportOpaqueMipTexture: typeof TSL.viewportOpaqueMipTexture;
export const viewportResolution: typeof TSL.viewportResolution;
export const viewportSafeUV: typeof TSL.viewportSafeUV;
export const viewportSharedTexture: typeof TSL.viewportSharedTexture;
export const viewportSize: typeof TSL.viewportSize;
export const viewportTexture: typeof TSL.viewportTexture;
export const viewportUV: typeof TSL.viewportUV;
export const wgsl: typeof TSL.wgsl;
export const wgslFn: typeof TSL.wgslFn;
export const workgroupArray: typeof TSL.workgroupArray;
export const workgroupBarrier: typeof TSL.workgroupBarrier;
export const workgroupId: typeof TSL.workgroupId;
export const workingToColorSpace: typeof TSL.workingToColorSpace;
export const xor: typeof TSL.xor;

export type { ProxiedObject } from "./nodes/TSL.js";


// @filename: /node_modules/three/src/Three.WebGPU.d.ts
export { default as Node } from "./nodes/core/Node.js";
export * from "./nodes/Nodes.js";


// @filename: /node_modules/three/src/nodes/TSL.d.ts
// constants
export * from "./core/constants.js";

// core
export * from "./core/AssignNode.js";
export * from "./core/AttributeNode.js";
export * from "./core/BypassNode.js";
export * from "./core/ContextNode.js";
export * from "./core/IndexNode.js";
export * from "./core/IsolateNode.js";
export * from "./core/MRTNode.js";
export * from "./core/OutputStructNode.js";
export * from "./core/OverrideContextNode.js";
export * from "./core/ParameterNode.js";
export * from "./core/PropertyNode.js";
export * from "./core/StackNode.js";
export * from "./core/StructNode.js";
export * from "./core/UniformGroupNode.js";
export * from "./core/UniformNode.js";
export * from "./core/VaryingNode.js";

// math
export * from "./math/BitcastNode.js";
export * from "./math/BitcountNode.js";
export * from "./math/Hash.js";
export * from "./math/MathUtils.js";
export * from "./math/PackFloatNode.js";
export * from "./math/TriNoise3D.js";
export * from "./math/UnpackFloatNode.js";

// utils
export * from "./utils/EquirectUV.js";
export * from "./utils/EventNode.js";
export * from "./utils/FunctionOverloadingNode.js";
export * from "./utils/LoopNode.js";
export * from "./utils/MatcapUV.js";
export * from "./utils/MaxMipLevelNode.js";
export * from "./utils/Oscillators.js";
export * from "./utils/Packing.js";
export * from "./utils/PostProcessingUtils.js";
export * from "./utils/ReflectorNode.js";
export * from "./utils/Remap.js";
export * from "./utils/RotateNode.js";
export * from "./utils/RTTNode.js";
export * from "./utils/SampleNode.js";
export * from "./utils/SpriteSheetUV.js";
export * from "./utils/SpriteUtils.js";
export * from "./utils/Timer.js";
export * from "./utils/TriplanarTextures.js";
export * from "./utils/UVUtils.js";
export * from "./utils/ViewportUtils.js";

// three.js shading language
export * from "./tsl/TSLBase.js";

// accessors
export * from "./accessors/AccessorsUtils.js";
export * from "./accessors/Arrays.js";
export * from "./accessors/Batch.js";
export * from "./accessors/Bitangent.js";
export * from "./accessors/BufferAttributeNode.js";
export * from "./accessors/BufferNode.js";
export * from "./accessors/BuiltinNode.js";
export * from "./accessors/Camera.js";
export * from "./accessors/CubeTextureNode.js";
export * from "./accessors/Instance.js";
export * from "./accessors/MaterialNode.js";
export * from "./accessors/MaterialProperties.js";
export * from "./accessors/MaterialReferenceNode.js";
export * from "./accessors/ModelNode.js";
export * from "./accessors/ModelViewProjectionNode.js";
export * from "./accessors/Morph.js";
export * from "./accessors/Normal.js";
export * from "./accessors/Object3DNode.js";
export * from "./accessors/PointUVNode.js";
export * from "./accessors/Position.js";
export * from "./accessors/ReferenceNode.js";
export * from "./accessors/ReflectVector.js";
export * from "./accessors/RendererReferenceNode.js";
export * from "./accessors/SceneProperties.js";
export * from "./accessors/Skinning.js";
export * from "./accessors/StorageBufferNode.js";
export * from "./accessors/StorageTexture3DNode.js";
export * from "./accessors/StorageTextureNode.js";
export * from "./accessors/Tangent.js";
export * from "./accessors/Texture3DNode.js";
export * from "./accessors/TextureBicubic.js";
export * from "./accessors/TextureNode.js";
export * from "./accessors/TextureSizeNode.js";
export * from "./accessors/UniformArrayNode.js";
export * from "./accessors/UserDataNode.js";
export * from "./accessors/UV.js";
export * from "./accessors/VelocityNode.js";
export * from "./accessors/VertexColorNode.js";

// display
export * from "./display/BlendModes.js";
export * from "./display/BumpMapNode.js";
export * from "./display/ColorAdjustment.js";
export * from "./display/ColorSpaceNode.js";
export * from "./display/FrontFacingNode.js";
export * from "./display/NormalMapNode.js";
export * from "./display/PremultiplyAlphaFunctions.js";
export * from "./display/RenderOutputNode.js";
export * from "./display/ScreenNode.js";
export * from "./display/ToneMappingNode.js";
export * from "./display/ToonOutlinePassNode.js";
export * from "./display/ViewportDepthNode.js";
export * from "./display/ViewportDepthTextureNode.js";
export * from "./display/ViewportSharedTextureNode.js";
export * from "./display/ViewportTextureNode.js";

export * from "./display/PassNode.js";

export * from "./display/ColorSpaceFunctions.js";
export * from "./display/ToneMappingFunctions.js";

// code
export * from "./code/CodeNode.js";
export * from "./code/ExpressionNode.js";
export * from "./code/FunctionCallNode.js";
export * from "./code/FunctionNode.js";

// fog
export * from "./fog/Fog.js";

// geometry
export * from "./geometry/RangeNode.js";

// gpgpu
export * from "./gpgpu/AtomicFunctionNode.js";
export * from "./gpgpu/BarrierNode.js";
export * from "./gpgpu/ComputeBuiltinNode.js";
export * from "./gpgpu/ComputeNode.js";
export * from "./gpgpu/SubgroupFunctionNode.js";
export * from "./gpgpu/WorkgroupInfoNode.js";

// lighting
export * from "./accessors/Lights.js";
export * from "./lighting/LightingContextNode.js";
export * from "./lighting/LightsNode.js";
export * from "./lighting/PointLightNode.js";
export * from "./lighting/PointShadowNode.js";
export * from "./lighting/ShadowBaseNode.js";
export * from "./lighting/ShadowFilterNode.js";
export * from "./lighting/ShadowNode.js";

// pmrem
export * from "./pmrem/PMREMNode.js";
export * from "./pmrem/PMREMUtils.js";

// procedural
export * from "./procedural/Checker.js";

// shapes
export * from "./shapes/Shapes.js";

// materialX
export * from "./materialx/MaterialXNodes.js";

// functions
export { default as BRDF_GGX } from "./functions/BSDF/BRDF_GGX.js";
export { default as BRDF_Lambert } from "./functions/BSDF/BRDF_Lambert.js";
export { default as D_GGX } from "./functions/BSDF/D_GGX.js";
export { default as DFGLUT } from "./functions/BSDF/DFGLUT.js";
export { default as F_Schlick } from "./functions/BSDF/F_Schlick.js";
export { default as Schlick_to_F0 } from "./functions/BSDF/Schlick_to_F0.js";
export { default as V_GGX_SmithCorrelated } from "./functions/BSDF/V_GGX_SmithCorrelated.js";

export * from "./lighting/LightUtils.js";

export { default as getGeometryRoughness } from "./functions/material/getGeometryRoughness.js";
export { default as getParallaxCorrectNormal } from "./functions/material/getParallaxCorrectNormal.js";
export { default as getRoughness } from "./functions/material/getRoughness.js";
export { default as getShIrradianceAt } from "./functions/material/getShIrradianceAt.js";


// @filename: /node_modules/three/src/nodes/core/ArrayNode.d.ts
import Node from "./Node.js";
import TempNode from "./TempNode.js";

export interface ArrayNodeInterface<TNodeType> {
    count: number;
    values: Node<TNodeType>[] | null;
    readonly isArrayNode: true;
}

declare const ArrayNode: {
    new<TNodeType>(
        nodeType: TNodeType | null,
        count: number,
        values?: Node<TNodeType>[] | null,
    ): ArrayNode<TNodeType>;
};

type ArrayNode<TNodeType> = TempNode<TNodeType> & ArrayNodeInterface<TNodeType>;

export default ArrayNode;

interface ArrayFunction {
    <TNodeType>(values: Node<TNodeType>[]): ArrayNode<TNodeType>;
    <TNodeType extends string>(nodeType: TNodeType, count: number): ArrayNode<TNodeType>;
}

export const array: ArrayFunction;

declare module "./Node.js" {
    interface NodeExtensions<TNodeType> {
        toArray: (count: number) => ArrayNode<TNodeType>;
    }
}


// @filename: /node_modules/three/src/nodes/core/AssignNode.d.ts
import Node from "./Node.js";
import NodeBuilder from "./NodeBuilder.js";
import TempNode from "./TempNode.js";

declare class AssignNode extends TempNode {
    readonly isAssignNode: true;

    constructor(targetNode: Node, sourceNode: Node);

    needsSplitAssign(builder: NodeBuilder): boolean;
}

export default AssignNode;

export const assign: (targetNode: Node, sourceNode: Node | number) => AssignNode;


// @filename: /node_modules/three/src/nodes/core/AttributeNode.d.ts
import Node from "./Node.js";
import NodeBuilder from "./NodeBuilder.js";

interface AttributeNodeInterface {
    setAttributeName(attributeName: string): this;

    getAttributeName(builder: NodeBuilder): string;
}

declare const AttributeNode: {
    new<TNodeType>(attributeName: string, TNodeType?: string | null): AttributeNode<TNodeType>;
};

type AttributeNode<TNodeType = unknown> = Node<TNodeType> & AttributeNodeInterface;

export default AttributeNode;

export const attribute: <TNodeType>(
    name: string,
    nodeType?: TNodeType | null,
) => AttributeNode<TNodeType>;


// @filename: /node_modules/three/src/nodes/core/BypassNode.d.ts
import Node from "./Node.js";

export default class BypassNode extends Node {
    isBypassNode: true;
    outputNode: Node;
    callNode: Node;

    constructor(returnNode: Node, callNode: Node);
}

export const bypass: (returnNode: Node, callNode: Node) => BypassNode;

declare module "./Node.js" {
    interface NodeElements {
        bypass: (callNode: Node) => BypassNode;
    }
}


// @filename: /node_modules/three/src/nodes/core/ConstNode.d.ts
import InputNode from "./InputNode.js";
import NodeBuilder from "./NodeBuilder.js";

interface ConstNodeInterface {
    readonly isConstNode: true;

    generateConst(builder: NodeBuilder): string;
}

declare const ConstNode: {
    new<TNodeType, TValue>(value: TValue, nodeType?: TNodeType | null): ConstNode<TNodeType, TValue>;
};

type ConstNode<TNodeType, TValue> = InputNode<TNodeType, TValue> & ConstNodeInterface;

export default ConstNode;


// @filename: /node_modules/three/src/nodes/core/ContextNode.d.ts
import { Light } from "../../lights/Light.js";
import Node from "./Node.js";

declare class ContextNodeInterface<TNodeType> extends Node {
    readonly isContextNode: true;

    node: Node<TNodeType> | null;
    value: unknown;
}

declare const ContextNode: {
    new<TNodeType>(node?: Node<TNodeType> | null, value?: unknown): ContextNode<TNodeType>;
};

type ContextNode<TNodeType> = ContextNodeInterface<TNodeType> & Node<TNodeType>;

export default ContextNode;

interface ContextFunction {
    (value?: unknown): ContextNode<unknown>;
    <TNodeType>(node: Node<TNodeType>, value?: unknown): ContextNode<TNodeType>;
}

export const context: ContextFunction;

export const uniformFlow: <TNodeType>(node: Node<TNodeType>) => ContextNode<TNodeType>;

export const setName: <TNodeType>(node: Node<TNodeType>, label: string) => Node<TNodeType>;

export function builtinShadowContext(shadowNode: Node, light: Light, node?: Node | null): ContextNode<unknown>;

export function builtinAOContext(aoNode: Node, node?: Node | null): ContextNode<unknown>;

/**
 * @deprecated "label()" has been deprecated. Use "setName()" instead.
 */
export function label<TNodeType>(node: Node<TNodeType>, label: string): Node<TNodeType>;

declare module "./Node.js" {
    interface NodeExtensions<TNodeType> {
        context: (context?: unknown) => ContextNode<TNodeType>;

        /**
         * @deprecated "label()" has been deprecated. Use "setName()" instead.
         */
        label: (label: string) => Node<TNodeType>;

        uniformFlow: () => ContextNode<TNodeType>;

        setName: (label: string) => Node<TNodeType>;

        builtinShadowContext: (shadowNode: Node, light: Light) => ContextNode<TNodeType>;

        builtinAOContext: (aoValue: Node) => ContextNode<TNodeType>;
    }
}


// @filename: /node_modules/three/src/nodes/core/IndexNode.d.ts
import Node from "./Node.js";

export type IndexNodeScope =
    | typeof IndexNode.VERTEX
    | typeof IndexNode.INSTANCE
    | typeof IndexNode.SUBGROUP
    | typeof IndexNode.INVOCATION_LOCAL
    | typeof IndexNode.INVOCATION_SUBGROUP
    | typeof IndexNode.DRAW;

declare class IndexNode extends Node<"uint"> {
    scope: IndexNodeScope;

    readonly isInstanceNode: true;

    constructor(scope: IndexNodeScope);

    static VERTEX: "vertex";
    static INSTANCE: "instance";
    static SUBGROUP: "subgroup";
    static INVOCATION_LOCAL: "invocationLocal";
    static INVOCATION_SUBGROUP: "invocationSubgroup";
    static DRAW: "draw";
}

export default IndexNode;

export const vertexIndex: IndexNode;
export const instanceIndex: IndexNode;
export const subgroupIndex: IndexNode;
export const invocationSubgroupIndex: IndexNode;
export const invocationLocalIndex: IndexNode;
export const drawIndex: IndexNode;


// @filename: /node_modules/three/src/nodes/core/InputNode.d.ts
import Node from "./Node.js";
import NodeBuilder from "./NodeBuilder.js";

export type Precision = "low" | "medium" | "high";

interface InputNodeInterface<TValue> {
    isInputNode: true;
    value: TValue;
    precision: Precision | null;

    getInputType(builder: NodeBuilder): string | null;
    setPrecision(precision: Precision): this;
}

declare const InputNode: {
    new<TNodeType, TValue>(value: TValue, nodeType?: TNodeType | null): InputNode<TNodeType, TValue>;
};

type InputNode<TNodeType, TValue> = Node<TNodeType> & InputNodeInterface<TValue>;

export default InputNode;

export type { InputNode };


// @filename: /node_modules/three/src/nodes/core/InspectorNode.d.ts
import Node from "./Node.js";

declare class InspectorNode extends Node {
    constructor(node: Node, name?: string, callback?: (node: Node) => Node);
}

export default InspectorNode;

export function inspector<T extends Node>(node: T, name?: string, callback?: (node: T) => Node): T;

declare module "./Node.js" {
    interface NodeElements {
        toInspector: (name?: string, callback?: (node: this) => Node) => this;
    }
}


// @filename: /node_modules/three/src/nodes/core/IsolateNode.d.ts
import Node from "./Node.js";
import NodeCache from "./NodeCache.js";

declare class IsolateNode extends Node {
    node: Node;
    parent: boolean;

    readonly isIsolateNode: true;

    constructor(node: Node, parent?: boolean);
}

export default IsolateNode;

export const isolate: (node: Node) => IsolateNode;

/**
 * @deprecated "cache()" has been deprecated. Use "isolate()" instead.
 */
export const cache: (node: Node, cache?: NodeCache) => IsolateNode;

declare module "./Node.js" {
    interface NodeElements {
        /**
         * @deprecated "cache()" has been deprecated. Use "isolate()" instead.
         */
        cache: (cache?: NodeCache) => IsolateNode;

        isolate: () => IsolateNode;
    }
}


// @filename: /node_modules/three/src/nodes/core/MRTNode.d.ts
import BlendMode from "../../renderers/common/BlendMode.js";
import { Texture } from "../../textures/Texture.js";
import { Node } from "../Nodes.js";
import OutputStructNode from "./OutputStructNode.js";

export function getTextureIndex(textures: ReadonlyArray<Texture>, name: string): number;

declare class MRTNode extends OutputStructNode {
    outputNodes: { [name: string]: Node };

    blendModes: { [name: string]: BlendMode };

    readonly isMRTNode: true;

    constructor(outputNodes: { [name: string]: Node });

    setBlendMode(name: string, blend: BlendMode): this;

    getBlendMode(name: string): BlendMode;

    has(name: string): boolean;

    get: (name: string) => Node;

    merge(mrtNode: MRTNode): MRTNode;
}

export default MRTNode;

export const mrt: (outputNodes: { [name: string]: Node }) => MRTNode;


// @filename: /node_modules/three/src/nodes/core/Node.d.ts
import { EventDispatcher } from "../../core/EventDispatcher.js";
import { NodeUpdateType } from "./constants.js";
import NodeBuilder from "./NodeBuilder.js";
import NodeFrame from "./NodeFrame.js";

interface NodeJSONMeta {
    textures: {
        [key: string]: unknown;
    };
    images: {
        [key: string]: unknown;
    };
    nodes: {
        [key: string]: NodeJSONIntermediateOutputData;
    };
}

interface NodeJSONMetadata {
    version: number;
    type: "Node";
    generator: "Node.toJSON";
}

interface NodeJSONInputNodes {
    [property: string]:
        | string[]
        | {
            [index: string]: string | undefined;
        }
        | string
        | undefined;
}

export interface NodeJSONInputData {
    inputNodes?: NodeJSONInputNodes | undefined;
    meta: {
        textures: {
            [key: string]: unknown;
        };
        nodes: {
            [key: string]: Node;
        };
    };
}

export interface NodeJSONIntermediateOutputData {
    uuid: string;
    type: string | undefined;
    meta?: NodeJSONMeta | undefined;
    metadata?: NodeJSONMetadata;
    inputNodes?: NodeJSONInputNodes | undefined;
    textures?: unknown[];
    images?: unknown[];
    nodes?: NodeJSONIntermediateOutputData[];
}

interface NodeJSONOutputData {
    uuid: string;
    type: string | undefined;
    metadata?: NodeJSONMetadata;
    inputNodes?: NodeJSONInputNodes | undefined;
    textures?: unknown[];
    images?: unknown[];
    nodes?: NodeJSONOutputData[];
}

export interface NodeChild {
    property: string;
    index?: number | string;
    childNode: Node;
}

export interface NodeClassEventMap {
    dispose: {};
}

/**
 * Base class for all nodes.
 *
 * @augments EventDispatcher
 */
declare class NodeClass<TEventMap extends NodeClassEventMap = NodeClassEventMap> extends EventDispatcher<TEventMap> {
    static get type(): string;
    /**
     * Constructs a new node.
     *
     * @param {?string} nodeType - The node type.
     */
    constructor(nodeType?: string | null);
    /**
     * The node type. This represents the result type of the node (e.g. `float` or `vec3`).
     *
     * @type {?string}
     * @default null
     */
    nodeType: string | null;
    /**
     * The update type of the node's {@link Node#update} method. Possible values are listed in {@link NodeUpdateType}.
     *
     * @type {string}
     * @default 'none'
     */
    updateType: NodeUpdateType;
    /**
     * The update type of the node's {@link Node#updateBefore} method. Possible values are listed in {@link NodeUpdateType}.
     *
     * @type {string}
     * @default 'none'
     */
    updateBeforeType: NodeUpdateType;
    /**
     * The update type of the node's {@link Node#updateAfter} method. Possible values are listed in {@link NodeUpdateType}.
     *
     * @type {string}
     * @default 'none'
     */
    updateAfterType: NodeUpdateType;
    /**
     * The version of the node. The version automatically is increased when {@link Node#needsUpdate} is set to `true`.
     *
     * @type {number}
     * @readonly
     * @default 0
     */
    readonly version: number;
    /**
     * The name of the node.
     *
     * @type {string}
     * @default ''
     */
    name: string;
    /**
     * Whether this node is global or not. This property is relevant for the internal
     * node caching system. All nodes which should be declared just once should
     * set this flag to `true` (a typical example is {@link AttributeNode}).
     *
     * @type {boolean}
     * @default false
     */
    global: boolean;
    /**
     * Create a list of parents for this node during the build process.
     *
     * @type {boolean}
     * @default false
     */
    parents: boolean;
    /**
     * This flag can be used for type testing.
     *
     * @type {boolean}
     * @readonly
     * @default true
     */
    readonly isNode: boolean;
    /**
     * The unique ID of the node.
     *
     * @type {number}
     * @readonly
     */
    readonly id: number;
    /**
     * The stack trace of the node for debugging purposes.
     *
     * @type {?string}
     * @default null
     */
    stackTrace: string | null;
    /**
     * Set this property to `true` when the node should be regenerated.
     *
     * @type {boolean}
     * @default false
     * @param {boolean} value
     */
    set needsUpdate(value: boolean);
    /**
     * The UUID of the node.
     *
     * @type {string}
     * @readonly
     */
    get uuid(): string;
    /**
     * The type of the class. The value is usually the constructor name.
     *
     * @type {string}
     * @readonly
     */
    get type(): string;
    /**
     * Convenient method for defining {@link Node#update}.
     *
     * @param {Function} callback - The update method.
     * @param {string} updateType - The update type.
     * @return {Node} A reference to this node.
     */
    onUpdate(callback: (this: this, frame: NodeFrame) => unknown, updateType: NodeUpdateType): this;
    /**
     * The method can be implemented to update the node's internal state when it is used to render an object.
     * The {@link Node#updateType} property defines how often the update is executed.
     *
     * @abstract
     * @param {NodeFrame} frame - A reference to the current node frame.
     * @return {?boolean} An optional bool that indicates whether the implementation actually performed an update or not (e.g. due to caching).
     */
    update(frame: NodeFrame): boolean | undefined;
    /**
     * Convenient method for defining {@link Node#update}. Similar to {@link Node#onUpdate}, but
     * this method automatically sets the update type to `FRAME`.
     *
     * @param {Function} callback - The update method.
     * @return {Node} A reference to this node.
     */
    onFrameUpdate(callback: (this: this, frame: NodeFrame) => unknown): this;
    /**
     * Convenient method for defining {@link Node#update}. Similar to {@link Node#onUpdate}, but
     * this method automatically sets the update type to `RENDER`.
     *
     * @param {Function} callback - The update method.
     * @return {Node} A reference to this node.
     */
    onRenderUpdate(callback: (this: this, frame: NodeFrame) => unknown): this;
    /**
     * Convenient method for defining {@link Node#update}. Similar to {@link Node#onUpdate}, but
     * this method automatically sets the update type to `OBJECT`.
     *
     * @param {Function} callback - The update method.
     * @return {Node} A reference to this node.
     */
    onObjectUpdate(callback: (this: this, frame: NodeFrame) => unknown): this;
    /**
     * Convenient method for defining {@link Node#updateReference}.
     *
     * @param {Function} callback - The update method.
     * @return {Node} A reference to this node.
     */
    onReference(callback: (this: this, frame: NodeBuilder | NodeFrame) => unknown): this;
    /**
     * Nodes might refer to other objects like materials. This method allows to dynamically update the reference
     * to such objects based on a given state (e.g. the current node frame or builder).
     *
     * @param {any} state - This method can be invocated in different contexts so `state` can refer to any object type.
     * @return {any} The updated reference.
     */
    updateReference(state: NodeBuilder | NodeFrame): unknown;
    /**
     * By default this method returns the value of the {@link Node#global} flag. This method
     * can be overwritten in derived classes if an analytical way is required to determine the
     * global cache referring to the current shader-stage.
     *
     * @param {NodeBuilder} builder - The current node builder.
     * @return {boolean} Whether this node is global or not.
     */
    isGlobal(builder: NodeBuilder): boolean;
    /**
     * Generator function that can be used to iterate over the child nodes.
     *
     * @generator
     * @yields {Node} A child node.
     */
    getChildren(): Generator<Node<unknown>, void, unknown>;
    /**
     * Calling this method dispatches the `dispose` event. This event can be used
     * to register event listeners for clean up tasks.
     */
    dispose(): void;
    /**
     * Can be used to traverse through the node's hierarchy.
     *
     * @param {traverseCallback} callback - A callback that is executed per node.
     */
    traverse(callback: (node: Node) => void): void;
    /**
     * Returns the cache key for this node.
     *
     * @param {boolean} [force=false] - When set to `true`, a recomputation of the cache key is forced.
     * @param {Set<Node>} [ignores=null] - A set of nodes to ignore during the computation of the cache key.
     * @return {number} The cache key of the node.
     */
    getCacheKey(force?: boolean, ignores?: Set<Node>): number;
    /**
     * Generate a custom cache key for this node.
     *
     * @return {number} The cache key of the node.
     */
    customCacheKey(): number;
    /**
     * Returns the references to this node which is by default `this`.
     *
     * @return {Node} A reference to this node.
     */
    getScope(): this;
    /**
     * Returns the hash of the node which is used to identify the node. By default it's
     * the {@link Node#uuid} however derived node classes might have to overwrite this method
     * depending on their implementation.
     *
     * @param {NodeBuilder} builder - The current node builder.
     * @return {string} The hash.
     */
    getHash(builder: NodeBuilder): string;
    /**
     * Returns the update type of {@link Node#update}.
     *
     * @return {NodeUpdateType} The update type.
     */
    getUpdateType(): NodeUpdateType;
    /**
     * Returns the update type of {@link Node#updateBefore}.
     *
     * @return {NodeUpdateType} The update type.
     */
    getUpdateBeforeType(): NodeUpdateType;
    /**
     * Returns the update type of {@link Node#updateAfter}.
     *
     * @return {NodeUpdateType} The update type.
     */
    getUpdateAfterType(): NodeUpdateType;
    /**
     * Certain types are composed of multiple elements. For example a `vec3`
     * is composed of three `float` values. This method returns the type of
     * these elements.
     *
     * @param {NodeBuilder} builder - The current node builder.
     * @return {string} The type of the node.
     */
    getElementType(builder: NodeBuilder): string;
    /**
     * Returns the node member type for the given name.
     *
     * @param {NodeBuilder} builder - The current node builder.
     * @param {string} name - The name of the member.
     * @return {string} The type of the node.
     */
    getMemberType(builder: NodeBuilder, name: string): string;
    /**
     * Returns the node's type.
     *
     * @param {NodeBuilder} builder - The current node builder.
     * @param {string} [output=null] - The output of the node.
     * @return {string} The type of the node.
     */
    getNodeType(builder: NodeBuilder, output?: string): string;
    /**
     * Returns the node's type.
     *
     * @param {NodeBuilder} builder - The current node builder.
     * @param {string} [output=null] - The output of the node.
     * @return {string} The type of the node.
     */
    generateNodeType(builder: NodeBuilder, output?: string): string;
    /**
     * This method is used during the build process of a node and ensures
     * equal nodes are not built multiple times but just once. For example if
     * `attribute( 'uv' )` is used multiple times by the user, the build
     * process makes sure to process just the first node. It also handles
     * node overrides if an override context is set.
     *
     * @param {NodeBuilder} builder - The current node builder.
     * @return {Node} The shared node if possible. Otherwise `this` is returned.
     */
    getShared(builder: NodeBuilder): Node;
    /**
     * Returns the number of elements in the node array.
     *
     * @param {NodeBuilder} builder - The current node builder.
     * @return {?number} The number of elements in the node array.
     */
    getArrayCount(builder: NodeBuilder): number | null;
    /**
     * Represents the setup stage which is the first step of the build process, see {@link Node#build} method.
     * This method is often overwritten in derived modules to prepare the node which is used as a node's output/result.
     * If an output node is prepared, then it must be returned in the `return` statement of the derived module's setup function.
     *
     * @param {NodeBuilder} builder - The current node builder.
     * @return {?Node} The output node.
     */
    setup(builder: NodeBuilder): Node | null | undefined;
    /**
     * Represents the analyze stage which is the second step of the build process, see {@link Node#build} method.
     * This stage analyzes the node hierarchy and ensures descendent nodes are built.
     *
     * @param {NodeBuilder} builder - The current node builder.
     * @param {?Node} output - The target output node.
     */
    analyze(builder: NodeBuilder, output?: Node | null): void;
    /**
     * Represents the generate stage which is the third step of the build process, see {@link Node#build} method.
     * This state builds the output node and returns the resulting shader string.
     *
     * @param {NodeBuilder} builder - The current node builder.
     * @param {?string} [output] - Can be used to define the output type.
     * @return {?string} The generated shader string.
     */
    generate(builder: NodeBuilder, output?: string | null): string | null | undefined;
    /**
     * The method can be implemented to update the node's internal state before it is used to render an object.
     * The {@link Node#updateBeforeType} property defines how often the update is executed.
     *
     * @abstract
     * @param {NodeFrame} frame - A reference to the current node frame.
     * @return {?boolean} An optional bool that indicates whether the implementation actually performed an update or not (e.g. due to caching).
     */
    updateBefore(frame: NodeFrame): boolean | undefined;
    /**
     * The method can be implemented to update the node's internal state after it was used to render an object.
     * The {@link Node#updateAfterType} property defines how often the update is executed.
     *
     * @abstract
     * @param {NodeFrame} frame - A reference to the current node frame.
     * @return {?boolean} An optional bool that indicates whether the implementation actually performed an update or not (e.g. due to caching).
     */
    updateAfter(frame: NodeFrame): boolean | undefined;
    before(node: Node): this;
    /**
     * This method performs the build of a node. The behavior and return value depend on the current build stage:
     * - **setup**: Prepares the node and its children for the build process. This process can also create new nodes. Returns the node itself or a variant.
     * - **analyze**: Analyzes the node hierarchy for optimizations in the code generation stage. Returns `null`.
     * - **generate**: Generates the shader code for the node. Returns the generated shader string.
     *
     * @param {NodeBuilder} builder - The current node builder.
     * @param {?(string|Node)} [output=null] - Can be used to define the output type.
     * @return {?(Node|string)} The result of the build process, depending on the build stage.
     */
    build(builder: NodeBuilder, output?: (string | Node) | null): (Node | string) | null;
    /**
     * Returns the child nodes as a JSON object.
     *
     * @return {Generator<Object>} An iterable list of serialized child objects as JSON.
     */
    getSerializeChildren(): NodeChild[];
    /**
     * Serializes the node to JSON.
     *
     * @param {Object} json - The output JSON object.
     */
    serialize(json: NodeJSONIntermediateOutputData): void;
    /**
     * Deserializes the node from the given JSON.
     *
     * @param {Object} json - The JSON object.
     */
    deserialize(json: NodeJSONInputData): void;
    /**
     * Serializes the node into the three.js JSON Object/Scene format.
     *
     * @param {?Object} meta - An optional JSON object that already holds serialized data from other scene objects.
     * @return {Object} The serialized node.
     */
    toJSON(meta?: NodeJSONMeta | string): NodeJSONOutputData;
}

declare const Node: {
    new<TNodeType>(nodeType?: TNodeType | null): Node<TNodeType>;
    new(nodeType?: string | null): Node;
    /**
     * Enables or disables the automatic capturing of stack traces for nodes.
     *
     * @type {boolean}
     * @default false
     */
    captureStackTrace: boolean;
};

export interface NodeElements {
}

export interface NodeExtensions<TNodeType> {
}

export type NumType = "float" | "int" | "uint";
export type IntegerType = "int" | "uint";
export type NumOrBoolType = NumType | "bool";
export type FloatVecType = "vec2" | "vec3" | "vec4";
export type MatType = "mat2" | "mat3" | "mat4";

export interface FloatExtensions {
}

export interface IntExtensions {
}

export interface UintExtensions {
}

export interface BoolExtensions {
}

export interface NumExtensions<TNum extends NumType> {
}

export interface IntegerExtensions<TInteger extends IntegerType> {
}

export interface NumOrBoolExtensions<TNumOrBool extends NumOrBoolType> {
}

export interface NumVec2Extensions<TNum extends NumType> {
}

export interface NumVec3Extensions<TNum extends NumType> {
}

export interface NumVec4Extensions<TNum extends NumType> {
}

export interface NumOrBoolVec2Extensions<TNumOrBool extends NumOrBoolType> {
}

export interface NumOrBoolVec3Extensions<TNumOrBool extends NumOrBoolType> {
}

export interface NumOrBoolVec4Extensions<TNumOrBool extends NumOrBoolType> {
}

export interface Vec2Extensions {
}

export interface Vec3Extensions {
}

export interface Vec4Extensions {
}

export interface ColorExtensions {
}

export interface FloatVecExtensions<TVec extends FloatVecType> {
}

export interface FloatOrVecExtensions<TNodeType> {
}

export interface IntOrVecExtensions<TNodeType> {
}

export interface UintOrVecExtensions<TNodeType> {
}

export interface BoolOrVecExtensions<TNodeType> {
}

export interface Mat2Extensions {
}

export interface Mat3Extensions {
}

export interface Mat4Extensions {
}

export interface MatExtensions<TMat extends MatType> {
}

type Node<TNodeType = unknown> =
    & NodeClass
    & NodeElements
    & (unknown extends TNodeType ? {} : NodeExtensions<TNodeType>)
    & (TNodeType extends "float"
        ? NumOrBoolExtensions<"float"> & FloatExtensions & NumExtensions<"float"> & FloatOrVecExtensions<"float">
        : TNodeType extends "int" ?
                & NumOrBoolExtensions<"int">
                & IntExtensions
                & NumExtensions<"int">
                & IntegerExtensions<"int">
                & IntOrVecExtensions<"int">
        : TNodeType extends "uint" ?
                & NumOrBoolExtensions<"uint">
                & UintExtensions
                & NumExtensions<"uint">
                & IntegerExtensions<"uint">
                & IntOrVecExtensions<"uint">
        : TNodeType extends "bool" ? NumOrBoolExtensions<"bool"> & BoolExtensions & BoolOrVecExtensions<"bool">
        : TNodeType extends "vec2" ?
                & NumOrBoolVec2Extensions<"float">
                & Vec2Extensions
                & NumVec2Extensions<"float">
                & FloatVecExtensions<"vec2">
                & FloatOrVecExtensions<"vec2">
        : TNodeType extends "ivec2"
            ? NumOrBoolVec2Extensions<"int"> & NumVec2Extensions<"int"> & IntOrVecExtensions<"ivec2">
        : TNodeType extends "uvec2"
            ? NumOrBoolVec2Extensions<"uint"> & NumVec2Extensions<"uint"> & IntOrVecExtensions<"uvec2">
        : TNodeType extends "bvec2" ? NumOrBoolVec2Extensions<"bool"> & BoolOrVecExtensions<"bvec2">
        : TNodeType extends "vec3" ?
                & NumOrBoolVec3Extensions<"float">
                & Vec3Extensions
                & NumVec3Extensions<"float">
                & FloatVecExtensions<"vec3">
                & FloatOrVecExtensions<"vec3">
        : TNodeType extends "ivec3"
            ? NumOrBoolVec3Extensions<"int"> & NumVec3Extensions<"int"> & IntOrVecExtensions<"ivec3">
        : TNodeType extends "uvec3"
            ? NumOrBoolVec3Extensions<"uint"> & NumVec3Extensions<"uint"> & IntOrVecExtensions<"uvec3">
        : TNodeType extends "bvec3" ? NumOrBoolVec3Extensions<"bool"> & BoolOrVecExtensions<"bvec3">
        : TNodeType extends "vec4" ?
                & NumOrBoolVec4Extensions<"float">
                & Vec4Extensions
                & NumVec4Extensions<"float">
                & FloatVecExtensions<"vec4">
                & FloatOrVecExtensions<"vec4">
        : TNodeType extends "ivec4"
            ? NumOrBoolVec4Extensions<"int"> & NumVec4Extensions<"int"> & IntOrVecExtensions<"ivec4">
        : TNodeType extends "uvec4"
            ? NumOrBoolVec4Extensions<"uint"> & NumVec4Extensions<"uint"> & IntOrVecExtensions<"uvec4">
        : TNodeType extends "bvec4" ? NumOrBoolVec4Extensions<"bool"> & BoolOrVecExtensions<"bvec4">
        : TNodeType extends "color" ? ColorExtensions
        : TNodeType extends "mat2" ? Mat2Extensions & MatExtensions<"mat2">
        : TNodeType extends "mat3" ? Mat3Extensions & MatExtensions<"mat3">
        : TNodeType extends "mat4" ? Mat4Extensions & MatExtensions<"mat4">
        : {})
    & {
        __TypeScript_NODE_TYPE__: TNodeType;
    };

export default Node;

export {};


// @filename: /node_modules/three/src/nodes/core/NodeBuilder.d.ts
import { BufferGeometry } from "../../core/BufferGeometry.js";
import { Object3D } from "../../core/Object3D.js";
import { Material } from "../../materials/Material.js";
import Renderer from "../../renderers/common/Renderer.js";

export default abstract class NodeBuilder {
    object: Object3D;
    material: Material;
    geometry: BufferGeometry;
    renderer: Renderer;
    context: unknown;
}


// @filename: /node_modules/three/src/nodes/core/NodeCode.d.ts
export default class NodeCode {
    isNodeCode: true;
    constructor(name: string, type: string, code?: string);
}


// @filename: /node_modules/three/src/nodes/core/NodeError.d.ts
import StackTrace from "./StackTrace.js";

declare class NodeError extends Error {
    stackTrace: StackTrace | null;

    constructor(message?: string, stackTrace?: StackTrace | null);
}

export default NodeError;


// @filename: /node_modules/three/src/nodes/core/NodeFrame.d.ts
import { Camera } from "../../cameras/Camera.js";
import { Object3D } from "../../core/Object3D.js";
import { Material } from "../../materials/Material.js";
import Renderer from "../../renderers/common/Renderer.js";
import { Scene } from "../../scenes/Scene.js";
import Node from "./Node.js";

export default class NodeFrame {
    time: number;
    deltaTime: number;

    frameId: number;
    renderId: number;

    startTime: number | null;

    frameMap: WeakMap<Node, number>;
    frameBeforeMap: WeakMap<Node, number>;
    renderMap: WeakMap<Node, number>;
    renderBeforeMap: WeakMap<Node, number>;

    renderer: Renderer | null;
    material: Material | null;
    camera: Camera | null;
    object: Object3D | null;
    scene: Scene | null;

    constructor();

    updateBeforeNode(node: Node): void;

    updateNode(node: Node): void;
    update(): void;
}


// @filename: /node_modules/three/src/nodes/core/NodeFunction.d.ts
import NodeFunctionInput from "./NodeFunctionInput.js";

export default abstract class NodeFunction {
    isNodeFunction: true;
    type: string;
    inputs: NodeFunctionInput[];
    name: string;
    precision: string;

    constructor(type: string, inputs: NodeFunctionInput[], name?: string, precision?: string);

    abstract getCode(name?: string): string;
}


// @filename: /node_modules/three/src/nodes/core/NodeFunctionInput.d.ts
export default class NodeFunctionInput {
    isNodeFunctionInput: true;
    count: null | number;
    qualifier: string;
    isConst: boolean;
    constructor(type: string, name: string, count?: number, qualifier?: string, isConst?: boolean);
}


// @filename: /node_modules/three/src/nodes/core/NodeParser.d.ts
import NodeFunction from "./NodeFunction.js";

/**
 * Base class for node parsers. A derived parser must be implemented
 * for each supported native shader language.
 */
declare abstract class NodeParser {
    /**
     * The method parses the given native code an returns a node function.
     *
     * @abstract
     * @param {string} source - The native shader code.
     * @return {NodeFunction} A node function.
     */
    abstract parseFunction(source: string): NodeFunction;
}

export default NodeParser;


// @filename: /node_modules/three/src/nodes/core/NodeUtils.d.ts
import { Color } from "../../math/Color.js";
import { Matrix3 } from "../../math/Matrix3.js";
import { Matrix4 } from "../../math/Matrix4.js";
import { Vector2 } from "../../math/Vector2.js";
import { Vector3 } from "../../math/Vector3.js";
import { Vector4 } from "../../math/Vector4.js";
import Node from "./Node.js";

export const hashString: (str: string) => number;
export const hashArray: (array: number[]) => number;
export const hash: (...params: number[]) => number;

export function getTypeFromLength(length: number): string | undefined;

export function getLengthFromType(type: string): number | undefined;

export function getMemoryLengthFromType(type: string): number | undefined;

export function getAlignmentFromType(type: string): number | undefined;

export function getValueType(value: unknown): string | null;

export function getValueFromType(
    type: string,
    ...params: number[]
): Color | Vector2 | Vector3 | Vector4 | Matrix3 | Matrix4 | boolean | number | string | ArrayBufferLike | null;


// @filename: /node_modules/three/src/nodes/core/OutputStructNode.d.ts
import Node from "./Node.js";

export default class OutputStructNode extends Node {
    members: Node[];

    readonly isOutputStructNode: true;

    constructor(...members: Node[]);
}

export const outputStruct: (...members: Node[]) => OutputStructNode;


// @filename: /node_modules/three/src/nodes/core/ParameterNode.d.ts
import PropertyNode from "./PropertyNode.js";

interface ParameterNodeInterface {
    readonly isParameterNode: true;
}

declare const ParameterNode: {
    new<TNodeType>(
        nodeType: TNodeType,
        name?: string | null,
    ): ParameterNode<TNodeType>;
};

type ParameterNode<TNodeType> = PropertyNode<TNodeType> & ParameterNodeInterface;

export default ParameterNode;

export const parameter: <TNodeType>(type: TNodeType, name?: string | null) => ParameterNode<TNodeType>;


// @filename: /node_modules/three/src/nodes/core/StackNode.d.ts
import Node from "./Node.js";

declare class StackNode extends Node {
    isStackNode: true;
    nodes: Node[];
    outputNode: Node | null;

    constructor();

    If(boolNode: Node, method: () => void): this;

    ElseIf(boolNode: Node, method: () => void): this;

    Else(method: () => void): this;

    Switch(expression: Node): this;

    Case(...params: [...Node[], () => void]): this;

    Default(method: () => void): this;
}

export default StackNode;

export const stack: () => StackNode;


// @filename: /node_modules/three/src/nodes/core/StackTrace.d.ts
declare class StackTrace {
    readonly isStackTrace: true;

    stack: Array<{ fn: string; file: string; line: number; column: number }>;

    constructor(stackMessage?: Error | string | null);
}

export default StackTrace;


// @filename: /node_modules/three/src/nodes/core/StructNode.d.ts
import Node from "./Node.js";
import StructTypeNode, { MembersLayout } from "./StructTypeNode.js";

declare class StructNode extends Node {
    values: Node[];

    constructor(structLayoutNode: StructTypeNode, values: Node[]);
}

export default StructNode;

export interface Struct {
    (): StructNode;
    (values: Node[]): StructNode;
    (...values: Node[]): StructNode;
}

export const struct: (membersLayout: MembersLayout, name?: string | null) => Struct;


// @filename: /node_modules/three/src/nodes/core/StructType.d.ts
import { MemberLayout } from "./StructTypeNode.js";

declare class StructType {
    constructor(name: string, members: MemberLayout[]);
    name: string;
    members: MemberLayout[];
    output: boolean;
}

export default StructType;


// @filename: /node_modules/three/src/nodes/core/SubBuildNode.d.ts
import Node from "./Node.js";

declare class SubBuildNode extends Node {
    node: Node;
    name: string;

    readonly isSubBuildNode: true;

    constructor(node: Node, name: string, nodeType?: string | null);
}

export default SubBuildNode;

export const subBuild: (node: Node, name: string, type?: string | null) => SubBuildNode;


// @filename: /node_modules/three/src/nodes/core/TempNode.d.ts
import Node from "./Node.js";
import NodeBuilder from "./NodeBuilder.js";

interface TempNodeInterface {
    isTempNode: true;

    hasDependencies(builder: NodeBuilder): boolean;
}

declare const TempNode: {
    new<TNodeType = unknown>(type?: TNodeType | null): TempNode<TNodeType>;
};

type TempNode<TNodeType> = Node<TNodeType> & TempNodeInterface;

export default TempNode;


// @filename: /node_modules/three/src/nodes/core/UniformGroupNode.d.ts
import { NodeUpdateType } from "./constants.js";
import Node from "./Node.js";

export default class UniformGroupNode extends Node {
    name: string;
    version: number;

    shared: boolean;

    readonly isUniformGroup: true;

    constructor(name: string, shared?: boolean, order?: number, updateType?: NodeUpdateType | null);

    set needsUpdate(value: boolean);
}

export const uniformGroup: (name: string, order?: number, updateType?: NodeUpdateType | null) => UniformGroupNode;
export const sharedUniformGroup: (name: string, order?: number, updateType?: NodeUpdateType | null) => UniformGroupNode;

export const frameGroup: UniformGroupNode;
export const renderGroup: UniformGroupNode;
export const objectGroup: UniformGroupNode;


// @filename: /node_modules/three/src/nodes/tsl/TSLBase.d.ts
export * from "../accessors/BufferAttributeNode.js";
export * from "../code/ExpressionNode.js";
export * from "../code/FunctionCallNode.js";
export * from "../core/ArrayNode.js";
export * from "../core/AssignNode.js";
export * from "../core/BypassNode.js";
export * from "../core/ContextNode.js";
export * from "../core/InspectorNode.js";
export * from "../core/IsolateNode.js";
export * from "../core/PropertyNode.js";
export * from "../core/SubBuildNode.js";
export * from "../core/UniformNode.js";
export * from "../core/VarNode.js";
export * from "../core/VaryingNode.js";
export * from "../display/ColorSpaceNode.js";
export * from "../display/RenderOutputNode.js";
export * from "../display/ToneMappingNode.js";
export * from "../gpgpu/ComputeNode.js";
export * from "../math/ConditionalNode.js";
export * from "../math/MathNode.js";
export * from "../math/OperatorNode.js";
export * from "../utils/DebugNode.js";
export * from "../utils/Discard.js";
export * from "../utils/Remap.js";
export * from "./TSLCore.js";

/**
 * @deprecated
 */
export function addNodeElement(name: string): void;


// @filename: /node_modules/three/src/nodes/tsl/TSLCore.d.ts
import { Color } from "../../math/Color.js";
import { Matrix2 } from "../../math/Matrix2.js";
import { Matrix3 } from "../../math/Matrix3.js";
import { Matrix4 } from "../../math/Matrix4.js";
import { Vector2 } from "../../math/Vector2.js";
import { Vector3 } from "../../math/Vector3.js";
import { Vector4 } from "../../math/Vector4.js";
import ArrayNode from "../core/ArrayNode.js";
import ConstNode from "../core/ConstNode.js";
import Node, { NumOrBoolType } from "../core/Node.js";
import NodeBuilder from "../core/NodeBuilder.js";
import StackNode from "../core/StackNode.js";
import VarNode from "../core/VarNode.js";
import ArrayElementNode from "../utils/ArrayElementNode.js";
import ConvertNode from "../utils/ConvertNode.js";
import JoinNode from "../utils/JoinNode.js";

export function addMethodChaining(name: string, nodeElement: unknown): void;

declare module "../core/Node.js" {
    interface NodeElements {
        assign: (sourceNode: Node | number) => this;
        get: (value: string) => Node;
    }
}

interface NumOrBoolToJsType {
    float: number;
    int: number;
    uint: number;
    bool: boolean;
}
type NumOrBool<TNumOrBool extends NumOrBoolType> = Node<TNumOrBool> | NumOrBoolToJsType[TNumOrBool];

interface NumOrBoolToVec {
    float: "vec";
    int: "ivec";
    uint: "uvec";
    bool: "bvec";
}
type NumOrBoolToVec2<TNumOrBool extends NumOrBoolType> = `${NumOrBoolToVec[TNumOrBool]}2`;
type NumOrBoolToVec3<TNumOrBool extends NumOrBoolType> = `${NumOrBoolToVec[TNumOrBool]}3`;
type NumOrBoolToVec4<TNumOrBool extends NumOrBoolType> = `${NumOrBoolToVec[TNumOrBool]}4`;

interface Swizzle1From1<TNumOrBool extends NumOrBoolType> {
    get x(): Node<TNumOrBool>;
    set x(value: NumOrBool<TNumOrBool>);
    get r(): Node<TNumOrBool>;
    set r(value: NumOrBool<TNumOrBool>);
    get s(): Node<TNumOrBool>;
    set s(value: NumOrBool<TNumOrBool>);
}

interface Swizzle1From2<TNumOrBool extends NumOrBoolType> extends Swizzle1From1<TNumOrBool> {
    get y(): Node<TNumOrBool>;
    set y(value: NumOrBool<TNumOrBool>);
    get g(): Node<TNumOrBool>;
    set g(value: NumOrBool<TNumOrBool>);
    get t(): Node<TNumOrBool>;
    set t(value: NumOrBool<TNumOrBool>);
}

interface Swizzle1From3<TNumOrBool extends NumOrBoolType> extends Swizzle1From2<TNumOrBool> {
    get z(): Node<TNumOrBool>;
    set z(value: NumOrBool<TNumOrBool>);
    get b(): Node<TNumOrBool>;
    set b(value: NumOrBool<TNumOrBool>);
    get p(): Node<TNumOrBool>;
    set p(value: NumOrBool<TNumOrBool>);
}

interface Swizzle1From4<TNumOrBool extends NumOrBoolType> extends Swizzle1From3<TNumOrBool> {
    get w(): Node<TNumOrBool>;
    set w(value: NumOrBool<TNumOrBool>);
    get a(): Node<TNumOrBool>;
    set a(value: NumOrBool<TNumOrBool>);
    get q(): Node<TNumOrBool>;
    set q(value: NumOrBool<TNumOrBool>);
}

interface Swizzle2From1<TNumOrBool extends NumOrBoolType> {
    get xx(): Node<NumOrBoolToVec2<TNumOrBool>>;
    get rr(): Node<NumOrBoolToVec2<TNumOrBool>>;
    get ss(): Node<NumOrBoolToVec2<TNumOrBool>>;
}

interface Swizzle2From2<TNumOrBool extends NumOrBoolType> extends Swizzle2From1<TNumOrBool> {
    get xy(): Node<NumOrBoolToVec2<TNumOrBool>>;
    set xy(value: Node<NumOrBoolToVec2<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get rg(): Node<NumOrBoolToVec2<TNumOrBool>>;
    set rg(value: Node<NumOrBoolToVec2<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get st(): Node<NumOrBoolToVec2<TNumOrBool>>;
    set st(value: Node<NumOrBoolToVec2<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get yx(): Node<NumOrBoolToVec2<TNumOrBool>>;
    set yx(value: Node<NumOrBoolToVec2<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get gr(): Node<NumOrBoolToVec2<TNumOrBool>>;
    set gr(value: Node<NumOrBoolToVec2<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get ts(): Node<NumOrBoolToVec2<TNumOrBool>>;
    set ts(value: Node<NumOrBoolToVec2<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get yy(): Node<NumOrBoolToVec2<TNumOrBool>>;
    get gg(): Node<NumOrBoolToVec2<TNumOrBool>>;
    get tt(): Node<NumOrBoolToVec2<TNumOrBool>>;
}

interface Swizzle2From3<TNumOrBool extends NumOrBoolType> extends Swizzle2From2<TNumOrBool> {
    get xz(): Node<NumOrBoolToVec2<TNumOrBool>>;
    set xz(value: Node<NumOrBoolToVec2<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get rb(): Node<NumOrBoolToVec2<TNumOrBool>>;
    set rb(value: Node<NumOrBoolToVec2<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get sp(): Node<NumOrBoolToVec2<TNumOrBool>>;
    set sp(value: Node<NumOrBoolToVec2<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get yz(): Node<NumOrBoolToVec2<TNumOrBool>>;
    set yz(value: Node<NumOrBoolToVec2<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get gb(): Node<NumOrBoolToVec2<TNumOrBool>>;
    set gb(value: Node<NumOrBoolToVec2<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get tp(): Node<NumOrBoolToVec2<TNumOrBool>>;
    set tp(value: Node<NumOrBoolToVec2<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get zx(): Node<NumOrBoolToVec2<TNumOrBool>>;
    set zx(value: Node<NumOrBoolToVec2<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get br(): Node<NumOrBoolToVec2<TNumOrBool>>;
    set br(value: Node<NumOrBoolToVec2<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get ps(): Node<NumOrBoolToVec2<TNumOrBool>>;
    set ps(value: Node<NumOrBoolToVec2<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get zy(): Node<NumOrBoolToVec2<TNumOrBool>>;
    set zy(value: Node<NumOrBoolToVec2<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get bg(): Node<NumOrBoolToVec2<TNumOrBool>>;
    set bg(value: Node<NumOrBoolToVec2<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get pt(): Node<NumOrBoolToVec2<TNumOrBool>>;
    set pt(value: Node<NumOrBoolToVec2<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get zz(): Node<NumOrBoolToVec2<TNumOrBool>>;
    get bb(): Node<NumOrBoolToVec2<TNumOrBool>>;
    get pp(): Node<NumOrBoolToVec2<TNumOrBool>>;
}

interface Swizzle2From4<TNumOrBool extends NumOrBoolType> extends Swizzle2From3<TNumOrBool> {
    get xw(): Node<NumOrBoolToVec2<TNumOrBool>>;
    set xw(value: Node<NumOrBoolToVec2<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get ra(): Node<NumOrBoolToVec2<TNumOrBool>>;
    set ra(value: Node<NumOrBoolToVec2<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get sq(): Node<NumOrBoolToVec2<TNumOrBool>>;
    set sq(value: Node<NumOrBoolToVec2<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get yw(): Node<NumOrBoolToVec2<TNumOrBool>>;
    set yw(value: Node<NumOrBoolToVec2<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get ga(): Node<NumOrBoolToVec2<TNumOrBool>>;
    set ga(value: Node<NumOrBoolToVec2<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get tq(): Node<NumOrBoolToVec2<TNumOrBool>>;
    set tq(value: Node<NumOrBoolToVec2<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get zw(): Node<NumOrBoolToVec2<TNumOrBool>>;
    set zw(value: Node<NumOrBoolToVec2<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get ba(): Node<NumOrBoolToVec2<TNumOrBool>>;
    set ba(value: Node<NumOrBoolToVec2<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get pq(): Node<NumOrBoolToVec2<TNumOrBool>>;
    set pq(value: Node<NumOrBoolToVec2<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get wx(): Node<NumOrBoolToVec2<TNumOrBool>>;
    set wx(value: Node<NumOrBoolToVec2<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get ar(): Node<NumOrBoolToVec2<TNumOrBool>>;
    set ar(value: Node<NumOrBoolToVec2<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get qs(): Node<NumOrBoolToVec2<TNumOrBool>>;
    set qs(value: Node<NumOrBoolToVec2<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get wy(): Node<NumOrBoolToVec2<TNumOrBool>>;
    set wy(value: Node<NumOrBoolToVec2<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get ag(): Node<NumOrBoolToVec2<TNumOrBool>>;
    set ag(value: Node<NumOrBoolToVec2<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get qt(): Node<NumOrBoolToVec2<TNumOrBool>>;
    set qt(value: Node<NumOrBoolToVec2<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get wz(): Node<NumOrBoolToVec2<TNumOrBool>>;
    set wz(value: Node<NumOrBoolToVec2<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get ab(): Node<NumOrBoolToVec2<TNumOrBool>>;
    set ab(value: Node<NumOrBoolToVec2<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get qp(): Node<NumOrBoolToVec2<TNumOrBool>>;
    set qp(value: Node<NumOrBoolToVec2<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get ww(): Node<NumOrBoolToVec2<TNumOrBool>>;
    get aa(): Node<NumOrBoolToVec2<TNumOrBool>>;
    get qq(): Node<NumOrBoolToVec2<TNumOrBool>>;
}

interface Swizzle3From1<TNumOrBool extends NumOrBoolType> {
    get xxx(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get rrr(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get sss(): Node<NumOrBoolToVec3<TNumOrBool>>;
}

interface Swizzle3From2<TNumOrBool extends NumOrBoolType> extends Swizzle3From1<TNumOrBool> {
    get xxy(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get rrg(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get sst(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get xyx(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get rgr(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get sts(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get xyy(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get rgg(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get stt(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get yxx(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get grr(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get tss(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get yxy(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get grg(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get tst(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get yyx(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get ggr(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get tts(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get yyy(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get ggg(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get ttt(): Node<NumOrBoolToVec3<TNumOrBool>>;
}

interface Swizzle3From3<TNumOrBool extends NumOrBoolType> extends Swizzle3From2<TNumOrBool> {
    get xxz(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get rrb(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get ssp(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get xyz(): Node<NumOrBoolToVec3<TNumOrBool>>;
    set xyz(value: Node<NumOrBoolToVec3<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get rgb(): Node<NumOrBoolToVec3<TNumOrBool>>;
    set rgb(value: Node<NumOrBoolToVec3<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get stp(): Node<NumOrBoolToVec3<TNumOrBool>>;
    set stp(value: Node<NumOrBoolToVec3<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get xzx(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get rbr(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get sps(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get xzy(): Node<NumOrBoolToVec3<TNumOrBool>>;
    set xzy(value: Node<NumOrBoolToVec3<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get rbg(): Node<NumOrBoolToVec3<TNumOrBool>>;
    set rbg(value: Node<NumOrBoolToVec3<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get spt(): Node<NumOrBoolToVec3<TNumOrBool>>;
    set spt(value: Node<NumOrBoolToVec3<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get xzz(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get rbb(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get spp(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get yxz(): Node<NumOrBoolToVec3<TNumOrBool>>;
    set yxz(value: Node<NumOrBoolToVec3<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get grb(): Node<NumOrBoolToVec3<TNumOrBool>>;
    set grb(value: Node<NumOrBoolToVec3<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get tsp(): Node<NumOrBoolToVec3<TNumOrBool>>;
    set tsp(value: Node<NumOrBoolToVec3<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get yyz(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get ggb(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get ttp(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get yzx(): Node<NumOrBoolToVec3<TNumOrBool>>;
    set yzx(value: Node<NumOrBoolToVec3<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get gbr(): Node<NumOrBoolToVec3<TNumOrBool>>;
    set gbr(value: Node<NumOrBoolToVec3<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get tps(): Node<NumOrBoolToVec3<TNumOrBool>>;
    set tps(value: Node<NumOrBoolToVec3<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get yzy(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get gbg(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get tpt(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get yzz(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get gbb(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get tpp(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get zxx(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get brr(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get pss(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get zxy(): Node<NumOrBoolToVec3<TNumOrBool>>;
    set zxy(value: Node<NumOrBoolToVec3<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get brg(): Node<NumOrBoolToVec3<TNumOrBool>>;
    set brg(value: Node<NumOrBoolToVec3<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get pst(): Node<NumOrBoolToVec3<TNumOrBool>>;
    set pst(value: Node<NumOrBoolToVec3<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get zxz(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get brb(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get psp(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get zyx(): Node<NumOrBoolToVec3<TNumOrBool>>;
    set zyx(value: Node<NumOrBoolToVec3<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get bgr(): Node<NumOrBoolToVec3<TNumOrBool>>;
    set bgr(value: Node<NumOrBoolToVec3<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get pts(): Node<NumOrBoolToVec3<TNumOrBool>>;
    set pts(value: Node<NumOrBoolToVec3<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get zyy(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get bgg(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get ptt(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get zyz(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get bgb(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get ptp(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get zzx(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get bbr(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get pps(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get zzy(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get bbg(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get ppt(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get zzz(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get bbb(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get ppp(): Node<NumOrBoolToVec3<TNumOrBool>>;
}

interface Swizzle3From4<TNumOrBool extends NumOrBoolType> extends Swizzle3From3<TNumOrBool> {
    get xxw(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get rra(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get ssq(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get xyw(): Node<NumOrBoolToVec3<TNumOrBool>>;
    set xyw(value: Node<NumOrBoolToVec3<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get rga(): Node<NumOrBoolToVec3<TNumOrBool>>;
    set rga(value: Node<NumOrBoolToVec3<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get stq(): Node<NumOrBoolToVec3<TNumOrBool>>;
    set stq(value: Node<NumOrBoolToVec3<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get xzw(): Node<NumOrBoolToVec3<TNumOrBool>>;
    set xzw(value: Node<NumOrBoolToVec3<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get rba(): Node<NumOrBoolToVec3<TNumOrBool>>;
    set rba(value: Node<NumOrBoolToVec3<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get spq(): Node<NumOrBoolToVec3<TNumOrBool>>;
    set spq(value: Node<NumOrBoolToVec3<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get xwx(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get rar(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get sqs(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get xwy(): Node<NumOrBoolToVec3<TNumOrBool>>;
    set xwy(value: Node<NumOrBoolToVec3<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get rag(): Node<NumOrBoolToVec3<TNumOrBool>>;
    set rag(value: Node<NumOrBoolToVec3<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get sqt(): Node<NumOrBoolToVec3<TNumOrBool>>;
    set sqt(value: Node<NumOrBoolToVec3<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get xwz(): Node<NumOrBoolToVec3<TNumOrBool>>;
    set xwz(value: Node<NumOrBoolToVec3<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get rab(): Node<NumOrBoolToVec3<TNumOrBool>>;
    set rab(value: Node<NumOrBoolToVec3<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get sqp(): Node<NumOrBoolToVec3<TNumOrBool>>;
    set sqp(value: Node<NumOrBoolToVec3<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get xww(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get raa(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get sqq(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get yxw(): Node<NumOrBoolToVec3<TNumOrBool>>;
    set yxw(value: Node<NumOrBoolToVec3<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get gra(): Node<NumOrBoolToVec3<TNumOrBool>>;
    set gra(value: Node<NumOrBoolToVec3<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get tsq(): Node<NumOrBoolToVec3<TNumOrBool>>;
    set tsq(value: Node<NumOrBoolToVec3<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get yyw(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get gga(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get ttq(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get yzw(): Node<NumOrBoolToVec3<TNumOrBool>>;
    set yzw(value: Node<NumOrBoolToVec3<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get gba(): Node<NumOrBoolToVec3<TNumOrBool>>;
    set gba(value: Node<NumOrBoolToVec3<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get tpq(): Node<NumOrBoolToVec3<TNumOrBool>>;
    set tpq(value: Node<NumOrBoolToVec3<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get ywx(): Node<NumOrBoolToVec3<TNumOrBool>>;
    set ywx(value: Node<NumOrBoolToVec3<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get gar(): Node<NumOrBoolToVec3<TNumOrBool>>;
    set gar(value: Node<NumOrBoolToVec3<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get tqs(): Node<NumOrBoolToVec3<TNumOrBool>>;
    set tqs(value: Node<NumOrBoolToVec3<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get ywy(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get gag(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get tqt(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get ywz(): Node<NumOrBoolToVec3<TNumOrBool>>;
    set ywz(value: Node<NumOrBoolToVec3<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get gab(): Node<NumOrBoolToVec3<TNumOrBool>>;
    set gab(value: Node<NumOrBoolToVec3<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get tqp(): Node<NumOrBoolToVec3<TNumOrBool>>;
    set tqp(value: Node<NumOrBoolToVec3<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get yww(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get gaa(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get tqq(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get zxw(): Node<NumOrBoolToVec3<TNumOrBool>>;
    set zxw(value: Node<NumOrBoolToVec3<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get bra(): Node<NumOrBoolToVec3<TNumOrBool>>;
    set bra(value: Node<NumOrBoolToVec3<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get psq(): Node<NumOrBoolToVec3<TNumOrBool>>;
    set psq(value: Node<NumOrBoolToVec3<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get zyw(): Node<NumOrBoolToVec3<TNumOrBool>>;
    set zyw(value: Node<NumOrBoolToVec3<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get bga(): Node<NumOrBoolToVec3<TNumOrBool>>;
    set bga(value: Node<NumOrBoolToVec3<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get ptq(): Node<NumOrBoolToVec3<TNumOrBool>>;
    set ptq(value: Node<NumOrBoolToVec3<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get zzw(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get bba(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get ppq(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get zwx(): Node<NumOrBoolToVec3<TNumOrBool>>;
    set zwx(value: Node<NumOrBoolToVec3<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get bar(): Node<NumOrBoolToVec3<TNumOrBool>>;
    set bar(value: Node<NumOrBoolToVec3<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get pqs(): Node<NumOrBoolToVec3<TNumOrBool>>;
    set pqs(value: Node<NumOrBoolToVec3<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get zwy(): Node<NumOrBoolToVec3<TNumOrBool>>;
    set zwy(value: Node<NumOrBoolToVec3<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get bag(): Node<NumOrBoolToVec3<TNumOrBool>>;
    set bag(value: Node<NumOrBoolToVec3<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get pqt(): Node<NumOrBoolToVec3<TNumOrBool>>;
    set pqt(value: Node<NumOrBoolToVec3<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get zwz(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get bab(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get pqp(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get zww(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get baa(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get pqq(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get wxx(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get arr(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get qss(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get wxy(): Node<NumOrBoolToVec3<TNumOrBool>>;
    set wxy(value: Node<NumOrBoolToVec3<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get arg(): Node<NumOrBoolToVec3<TNumOrBool>>;
    set arg(value: Node<NumOrBoolToVec3<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get qst(): Node<NumOrBoolToVec3<TNumOrBool>>;
    set qst(value: Node<NumOrBoolToVec3<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get wxz(): Node<NumOrBoolToVec3<TNumOrBool>>;
    set wxz(value: Node<NumOrBoolToVec3<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get arb(): Node<NumOrBoolToVec3<TNumOrBool>>;
    set arb(value: Node<NumOrBoolToVec3<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get qsp(): Node<NumOrBoolToVec3<TNumOrBool>>;
    set qsp(value: Node<NumOrBoolToVec3<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get wxw(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get ara(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get qsq(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get wyx(): Node<NumOrBoolToVec3<TNumOrBool>>;
    set wyx(value: Node<NumOrBoolToVec3<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get agr(): Node<NumOrBoolToVec3<TNumOrBool>>;
    set agr(value: Node<NumOrBoolToVec3<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get qts(): Node<NumOrBoolToVec3<TNumOrBool>>;
    set qts(value: Node<NumOrBoolToVec3<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get wyy(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get agg(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get qtt(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get wyz(): Node<NumOrBoolToVec3<TNumOrBool>>;
    set wyz(value: Node<NumOrBoolToVec3<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get agb(): Node<NumOrBoolToVec3<TNumOrBool>>;
    set agb(value: Node<NumOrBoolToVec3<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get qtp(): Node<NumOrBoolToVec3<TNumOrBool>>;
    set qtp(value: Node<NumOrBoolToVec3<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get wyw(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get aga(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get qtq(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get wzx(): Node<NumOrBoolToVec3<TNumOrBool>>;
    set wzx(value: Node<NumOrBoolToVec3<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get abr(): Node<NumOrBoolToVec3<TNumOrBool>>;
    set abr(value: Node<NumOrBoolToVec3<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get qps(): Node<NumOrBoolToVec3<TNumOrBool>>;
    set qps(value: Node<NumOrBoolToVec3<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get wzy(): Node<NumOrBoolToVec3<TNumOrBool>>;
    set wzy(value: Node<NumOrBoolToVec3<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get abg(): Node<NumOrBoolToVec3<TNumOrBool>>;
    set abg(value: Node<NumOrBoolToVec3<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get qpt(): Node<NumOrBoolToVec3<TNumOrBool>>;
    set qpt(value: Node<NumOrBoolToVec3<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get wzz(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get abb(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get qpp(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get wzw(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get aba(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get qpq(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get wwx(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get aar(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get qqs(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get wwy(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get aag(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get qqt(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get wwz(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get aab(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get qqp(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get www(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get aaa(): Node<NumOrBoolToVec3<TNumOrBool>>;
    get qqq(): Node<NumOrBoolToVec3<TNumOrBool>>;
}

interface Swizzle4From1<TNumOrBool extends NumOrBoolType> {
    get xxxx(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get rrrr(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get ssss(): Node<NumOrBoolToVec4<TNumOrBool>>;
}

interface Swizzle4From2<TNumOrBool extends NumOrBoolType> extends Swizzle4From1<TNumOrBool> {
    get xxxy(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get rrrg(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get ssst(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get xxyx(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get rrgr(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get ssts(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get xxyy(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get rrgg(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get sstt(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get xyxx(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get rgrr(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get stss(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get xyxy(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get rgrg(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get stst(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get xyyx(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get rggr(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get stts(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get xyyy(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get rggg(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get sttt(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get yxxx(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get grrr(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get tsss(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get yxxy(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get grrg(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get tsst(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get yxyx(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get grgr(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get tsts(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get yxyy(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get grgg(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get tstt(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get yyxx(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get ggrr(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get ttss(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get yyxy(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get ggrg(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get ttst(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get yyyx(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get gggr(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get ttts(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get yyyy(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get gggg(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get tttt(): Node<NumOrBoolToVec4<TNumOrBool>>;
}

interface Swizzle4From3<TNumOrBool extends NumOrBoolType> extends Swizzle4From2<TNumOrBool> {
    get xxxz(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get rrrb(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get sssp(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get xxyz(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get rrgb(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get sstp(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get xxzx(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get rrbr(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get ssps(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get xxzy(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get rrbg(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get sspt(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get xxzz(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get rrbb(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get sspp(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get xyxz(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get rgrb(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get stsp(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get xyyz(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get rggb(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get sttp(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get xyzx(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get rgbr(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get stps(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get xyzy(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get rgbg(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get stpt(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get xyzz(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get rgbb(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get stpp(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get xzxx(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get rbrr(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get spss(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get xzxy(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get rbrg(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get spst(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get xzxz(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get rbrb(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get spsp(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get xzyx(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get rbgr(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get spts(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get xzyy(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get rbgg(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get sptt(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get xzyz(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get rbgb(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get sptp(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get xzzx(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get rbbr(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get spps(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get xzzy(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get rbbg(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get sppt(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get xzzz(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get rbbb(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get sppp(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get yxxz(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get grrb(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get tssp(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get yxyz(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get grgb(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get tstp(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get yxzx(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get grbr(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get tsps(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get yxzy(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get grbg(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get tspt(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get yxzz(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get grbb(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get tspp(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get yyxz(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get ggrb(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get ttsp(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get yyyz(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get gggb(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get tttp(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get yyzx(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get ggbr(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get ttps(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get yyzy(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get ggbg(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get ttpt(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get yyzz(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get ggbb(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get ttpp(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get yzxx(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get gbrr(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get tpss(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get yzxy(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get gbrg(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get tpst(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get yzxz(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get gbrb(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get tpsp(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get yzyx(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get gbgr(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get tpts(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get yzyy(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get gbgg(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get tptt(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get yzyz(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get gbgb(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get tptp(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get yzzx(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get gbbr(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get tpps(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get yzzy(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get gbbg(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get tppt(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get yzzz(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get gbbb(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get tppp(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get zxxx(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get brrr(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get psss(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get zxxy(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get brrg(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get psst(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get zxxz(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get brrb(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get pssp(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get zxyx(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get brgr(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get psts(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get zxyy(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get brgg(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get pstt(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get zxyz(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get brgb(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get pstp(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get zxzx(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get brbr(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get psps(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get zxzy(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get brbg(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get pspt(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get zxzz(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get brbb(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get pspp(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get zyxx(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get bgrr(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get ptss(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get zyxy(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get bgrg(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get ptst(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get zyxz(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get bgrb(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get ptsp(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get zyyx(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get bggr(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get ptts(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get zyyy(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get bggg(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get pttt(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get zyyz(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get bggb(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get pttp(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get zyzx(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get bgbr(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get ptps(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get zyzy(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get bgbg(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get ptpt(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get zyzz(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get bgbb(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get ptpp(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get zzxx(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get bbrr(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get ppss(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get zzxy(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get bbrg(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get ppst(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get zzxz(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get bbrb(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get ppsp(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get zzyx(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get bbgr(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get ppts(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get zzyy(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get bbgg(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get pptt(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get zzyz(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get bbgb(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get pptp(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get zzzx(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get bbbr(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get ppps(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get zzzy(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get bbbg(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get pppt(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get zzzz(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get bbbb(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get pppp(): Node<NumOrBoolToVec4<TNumOrBool>>;
}

interface Swizzle4From4<TNumOrBool extends NumOrBoolType> extends Swizzle4From3<TNumOrBool> {
    get xxxw(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get rrra(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get sssq(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get xxyw(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get rrga(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get sstq(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get xxzw(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get rrba(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get sspq(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get xxwx(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get rrar(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get ssqs(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get xxwy(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get rrag(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get ssqt(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get xxwz(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get rrab(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get ssqp(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get xxww(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get rraa(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get ssqq(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get xyxw(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get rgra(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get stsq(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get xyyw(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get rgga(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get sttq(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get xyzw(): Node<NumOrBoolToVec4<TNumOrBool>>;
    set xyzw(value: Node<NumOrBoolToVec4<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get rgba(): Node<NumOrBoolToVec4<TNumOrBool>>;
    set rgba(value: Node<NumOrBoolToVec4<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get stpq(): Node<NumOrBoolToVec4<TNumOrBool>>;
    set stpq(value: Node<NumOrBoolToVec4<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get xywx(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get rgar(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get stqs(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get xywy(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get rgag(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get stqt(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get xywz(): Node<NumOrBoolToVec4<TNumOrBool>>;
    set xywz(value: Node<NumOrBoolToVec4<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get rgab(): Node<NumOrBoolToVec4<TNumOrBool>>;
    set rgab(value: Node<NumOrBoolToVec4<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get stqp(): Node<NumOrBoolToVec4<TNumOrBool>>;
    set stqp(value: Node<NumOrBoolToVec4<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get xyww(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get rgaa(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get stqq(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get xzxw(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get rbra(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get spsq(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get xzyw(): Node<NumOrBoolToVec4<TNumOrBool>>;
    set xzyw(value: Node<NumOrBoolToVec4<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get rbga(): Node<NumOrBoolToVec4<TNumOrBool>>;
    set rbga(value: Node<NumOrBoolToVec4<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get sptq(): Node<NumOrBoolToVec4<TNumOrBool>>;
    set sptq(value: Node<NumOrBoolToVec4<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get xzzw(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get rbba(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get sppq(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get xzwx(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get rbar(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get spqs(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get xzwy(): Node<NumOrBoolToVec4<TNumOrBool>>;
    set xzwy(value: Node<NumOrBoolToVec4<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get rbag(): Node<NumOrBoolToVec4<TNumOrBool>>;
    set rbag(value: Node<NumOrBoolToVec4<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get spqt(): Node<NumOrBoolToVec4<TNumOrBool>>;
    set spqt(value: Node<NumOrBoolToVec4<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get xzwz(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get rbab(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get spqp(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get xzww(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get rbaa(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get spqq(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get xwxx(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get rarr(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get sqss(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get xwxy(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get rarg(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get sqst(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get xwxz(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get rarb(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get sqsp(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get xwxw(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get rara(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get sqsq(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get xwyx(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get ragr(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get sqts(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get xwyy(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get ragg(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get sqtt(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get xwyz(): Node<NumOrBoolToVec4<TNumOrBool>>;
    set xwyz(value: Node<NumOrBoolToVec4<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get ragb(): Node<NumOrBoolToVec4<TNumOrBool>>;
    set ragb(value: Node<NumOrBoolToVec4<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get sqtp(): Node<NumOrBoolToVec4<TNumOrBool>>;
    set sqtp(value: Node<NumOrBoolToVec4<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get xwyw(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get raga(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get sqtq(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get xwzx(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get rabr(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get sqps(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get xwzy(): Node<NumOrBoolToVec4<TNumOrBool>>;
    set xwzy(value: Node<NumOrBoolToVec4<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get rabg(): Node<NumOrBoolToVec4<TNumOrBool>>;
    set rabg(value: Node<NumOrBoolToVec4<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get sqpt(): Node<NumOrBoolToVec4<TNumOrBool>>;
    set sqpt(value: Node<NumOrBoolToVec4<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get xwzz(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get rabb(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get sqpp(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get xwzw(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get raba(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get sqpq(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get xwwx(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get raar(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get sqqs(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get xwwy(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get raag(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get sqqt(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get xwwz(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get raab(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get sqqp(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get xwww(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get raaa(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get sqqq(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get yxxw(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get grra(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get tssq(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get yxyw(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get grga(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get tstq(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get yxzw(): Node<NumOrBoolToVec4<TNumOrBool>>;
    set yxzw(value: Node<NumOrBoolToVec4<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get grba(): Node<NumOrBoolToVec4<TNumOrBool>>;
    set grba(value: Node<NumOrBoolToVec4<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get tspq(): Node<NumOrBoolToVec4<TNumOrBool>>;
    set tspq(value: Node<NumOrBoolToVec4<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get yxwx(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get grar(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get tsqs(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get yxwy(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get grag(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get tsqt(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get yxwz(): Node<NumOrBoolToVec4<TNumOrBool>>;
    set yxwz(value: Node<NumOrBoolToVec4<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get grab(): Node<NumOrBoolToVec4<TNumOrBool>>;
    set grab(value: Node<NumOrBoolToVec4<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get tsqp(): Node<NumOrBoolToVec4<TNumOrBool>>;
    set tsqp(value: Node<NumOrBoolToVec4<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get yxww(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get graa(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get tsqq(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get yyxw(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get ggra(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get ttsq(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get yyyw(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get ggga(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get tttq(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get yyzw(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get ggba(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get ttpq(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get yywx(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get ggar(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get ttqs(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get yywy(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get ggag(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get ttqt(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get yywz(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get ggab(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get ttqp(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get yyww(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get ggaa(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get ttqq(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get yzxw(): Node<NumOrBoolToVec4<TNumOrBool>>;
    set yzxw(value: Node<NumOrBoolToVec4<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get gbra(): Node<NumOrBoolToVec4<TNumOrBool>>;
    set gbra(value: Node<NumOrBoolToVec4<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get tpsq(): Node<NumOrBoolToVec4<TNumOrBool>>;
    set tpsq(value: Node<NumOrBoolToVec4<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get yzyw(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get gbga(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get tptq(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get yzzw(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get gbba(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get tppq(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get yzwx(): Node<NumOrBoolToVec4<TNumOrBool>>;
    set yzwx(value: Node<NumOrBoolToVec4<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get gbar(): Node<NumOrBoolToVec4<TNumOrBool>>;
    set gbar(value: Node<NumOrBoolToVec4<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get tpqs(): Node<NumOrBoolToVec4<TNumOrBool>>;
    set tpqs(value: Node<NumOrBoolToVec4<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get yzwy(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get gbag(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get tpqt(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get yzwz(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get gbab(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get tpqp(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get yzww(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get gbaa(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get tpqq(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get ywxx(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get garr(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get tqss(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get ywxy(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get garg(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get tqst(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get ywxz(): Node<NumOrBoolToVec4<TNumOrBool>>;
    set ywxz(value: Node<NumOrBoolToVec4<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get garb(): Node<NumOrBoolToVec4<TNumOrBool>>;
    set garb(value: Node<NumOrBoolToVec4<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get tqsp(): Node<NumOrBoolToVec4<TNumOrBool>>;
    set tqsp(value: Node<NumOrBoolToVec4<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get ywxw(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get gara(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get tqsq(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get ywyx(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get gagr(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get tqts(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get ywyy(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get gagg(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get tqtt(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get ywyz(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get gagb(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get tqtp(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get ywyw(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get gaga(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get tqtq(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get ywzx(): Node<NumOrBoolToVec4<TNumOrBool>>;
    set ywzx(value: Node<NumOrBoolToVec4<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get gabr(): Node<NumOrBoolToVec4<TNumOrBool>>;
    set gabr(value: Node<NumOrBoolToVec4<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get tqps(): Node<NumOrBoolToVec4<TNumOrBool>>;
    set tqps(value: Node<NumOrBoolToVec4<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get ywzy(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get gabg(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get tqpt(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get ywzz(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get gabb(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get tqpp(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get ywzw(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get gaba(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get tqpq(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get ywwx(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get gaar(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get tqqs(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get ywwy(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get gaag(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get tqqt(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get ywwz(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get gaab(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get tqqp(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get ywww(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get gaaa(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get tqqq(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get zxxw(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get brra(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get pssq(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get zxyw(): Node<NumOrBoolToVec4<TNumOrBool>>;
    set zxyw(value: Node<NumOrBoolToVec4<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get brga(): Node<NumOrBoolToVec4<TNumOrBool>>;
    set brga(value: Node<NumOrBoolToVec4<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get pstq(): Node<NumOrBoolToVec4<TNumOrBool>>;
    set pstq(value: Node<NumOrBoolToVec4<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get zxzw(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get brba(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get pspq(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get zxwx(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get brar(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get psqs(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get zxwy(): Node<NumOrBoolToVec4<TNumOrBool>>;
    set zxwy(value: Node<NumOrBoolToVec4<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get brag(): Node<NumOrBoolToVec4<TNumOrBool>>;
    set brag(value: Node<NumOrBoolToVec4<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get psqt(): Node<NumOrBoolToVec4<TNumOrBool>>;
    set psqt(value: Node<NumOrBoolToVec4<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get zxwz(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get brab(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get psqp(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get zxww(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get braa(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get psqq(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get zyxw(): Node<NumOrBoolToVec4<TNumOrBool>>;
    set zyxw(value: Node<NumOrBoolToVec4<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get bgra(): Node<NumOrBoolToVec4<TNumOrBool>>;
    set bgra(value: Node<NumOrBoolToVec4<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get ptsq(): Node<NumOrBoolToVec4<TNumOrBool>>;
    set ptsq(value: Node<NumOrBoolToVec4<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get zyyw(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get bgga(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get pttq(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get zyzw(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get bgba(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get ptpq(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get zywx(): Node<NumOrBoolToVec4<TNumOrBool>>;
    set zywx(value: Node<NumOrBoolToVec4<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get bgar(): Node<NumOrBoolToVec4<TNumOrBool>>;
    set bgar(value: Node<NumOrBoolToVec4<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get ptqs(): Node<NumOrBoolToVec4<TNumOrBool>>;
    set ptqs(value: Node<NumOrBoolToVec4<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get zywy(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get bgag(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get ptqt(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get zywz(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get bgab(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get ptqp(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get zyww(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get bgaa(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get ptqq(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get zzxw(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get bbra(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get ppsq(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get zzyw(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get bbga(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get pptq(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get zzzw(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get bbba(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get pppq(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get zzwx(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get bbar(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get ppqs(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get zzwy(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get bbag(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get ppqt(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get zzwz(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get bbab(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get ppqp(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get zzww(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get bbaa(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get ppqq(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get zwxx(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get barr(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get pqss(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get zwxy(): Node<NumOrBoolToVec4<TNumOrBool>>;
    set zwxy(value: Node<NumOrBoolToVec4<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get barg(): Node<NumOrBoolToVec4<TNumOrBool>>;
    set barg(value: Node<NumOrBoolToVec4<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get pqst(): Node<NumOrBoolToVec4<TNumOrBool>>;
    set pqst(value: Node<NumOrBoolToVec4<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get zwxz(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get barb(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get pqsp(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get zwxw(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get bara(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get pqsq(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get zwyx(): Node<NumOrBoolToVec4<TNumOrBool>>;
    set zwyx(value: Node<NumOrBoolToVec4<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get bagr(): Node<NumOrBoolToVec4<TNumOrBool>>;
    set bagr(value: Node<NumOrBoolToVec4<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get pqts(): Node<NumOrBoolToVec4<TNumOrBool>>;
    set pqts(value: Node<NumOrBoolToVec4<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get zwyy(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get bagg(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get pqtt(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get zwyz(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get bagb(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get pqtp(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get zwyw(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get baga(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get pqtq(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get zwzx(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get babr(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get pqps(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get zwzy(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get babg(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get pqpt(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get zwzz(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get babb(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get pqpp(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get zwzw(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get baba(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get pqpq(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get zwwx(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get baar(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get pqqs(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get zwwy(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get baag(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get pqqt(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get zwwz(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get baab(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get pqqp(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get zwww(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get baaa(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get pqqq(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get wxxx(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get arrr(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get qsss(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get wxxy(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get arrg(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get qsst(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get wxxz(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get arrb(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get qssp(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get wxxw(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get arra(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get qssq(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get wxyx(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get argr(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get qsts(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get wxyy(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get argg(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get qstt(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get wxyz(): Node<NumOrBoolToVec4<TNumOrBool>>;
    set wxyz(value: Node<NumOrBoolToVec4<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get argb(): Node<NumOrBoolToVec4<TNumOrBool>>;
    set argb(value: Node<NumOrBoolToVec4<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get qstp(): Node<NumOrBoolToVec4<TNumOrBool>>;
    set qstp(value: Node<NumOrBoolToVec4<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get wxyw(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get arga(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get qstq(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get wxzx(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get arbr(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get qsps(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get wxzy(): Node<NumOrBoolToVec4<TNumOrBool>>;
    set wxzy(value: Node<NumOrBoolToVec4<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get arbg(): Node<NumOrBoolToVec4<TNumOrBool>>;
    set arbg(value: Node<NumOrBoolToVec4<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get qspt(): Node<NumOrBoolToVec4<TNumOrBool>>;
    set qspt(value: Node<NumOrBoolToVec4<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get wxzz(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get arbb(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get qspp(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get wxzw(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get arba(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get qspq(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get wxwx(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get arar(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get qsqs(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get wxwy(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get arag(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get qsqt(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get wxwz(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get arab(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get qsqp(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get wxww(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get araa(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get qsqq(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get wyxx(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get agrr(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get qtss(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get wyxy(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get agrg(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get qtst(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get wyxz(): Node<NumOrBoolToVec4<TNumOrBool>>;
    set wyxz(value: Node<NumOrBoolToVec4<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get agrb(): Node<NumOrBoolToVec4<TNumOrBool>>;
    set agrb(value: Node<NumOrBoolToVec4<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get qtsp(): Node<NumOrBoolToVec4<TNumOrBool>>;
    set qtsp(value: Node<NumOrBoolToVec4<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get wyxw(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get agra(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get qtsq(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get wyyx(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get aggr(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get qtts(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get wyyy(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get aggg(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get qttt(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get wyyz(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get aggb(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get qttp(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get wyyw(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get agga(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get qttq(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get wyzx(): Node<NumOrBoolToVec4<TNumOrBool>>;
    set wyzx(value: Node<NumOrBoolToVec4<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get agbr(): Node<NumOrBoolToVec4<TNumOrBool>>;
    set agbr(value: Node<NumOrBoolToVec4<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get qtps(): Node<NumOrBoolToVec4<TNumOrBool>>;
    set qtps(value: Node<NumOrBoolToVec4<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get wyzy(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get agbg(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get qtpt(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get wyzz(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get agbb(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get qtpp(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get wyzw(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get agba(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get qtpq(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get wywx(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get agar(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get qtqs(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get wywy(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get agag(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get qtqt(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get wywz(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get agab(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get qtqp(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get wyww(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get agaa(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get qtqq(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get wzxx(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get abrr(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get qpss(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get wzxy(): Node<NumOrBoolToVec4<TNumOrBool>>;
    set wzxy(value: Node<NumOrBoolToVec4<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get abrg(): Node<NumOrBoolToVec4<TNumOrBool>>;
    set abrg(value: Node<NumOrBoolToVec4<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get qpst(): Node<NumOrBoolToVec4<TNumOrBool>>;
    set qpst(value: Node<NumOrBoolToVec4<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get wzxz(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get abrb(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get qpsp(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get wzxw(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get abra(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get qpsq(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get wzyx(): Node<NumOrBoolToVec4<TNumOrBool>>;
    set wzyx(value: Node<NumOrBoolToVec4<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get abgr(): Node<NumOrBoolToVec4<TNumOrBool>>;
    set abgr(value: Node<NumOrBoolToVec4<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get qpts(): Node<NumOrBoolToVec4<TNumOrBool>>;
    set qpts(value: Node<NumOrBoolToVec4<TNumOrBool>> | NumOrBool<TNumOrBool>);
    get wzyy(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get abgg(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get qptt(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get wzyz(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get abgb(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get qptp(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get wzyw(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get abga(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get qptq(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get wzzx(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get abbr(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get qpps(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get wzzy(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get abbg(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get qppt(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get wzzz(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get abbb(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get qppp(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get wzzw(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get abba(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get qppq(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get wzwx(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get abar(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get qpqs(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get wzwy(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get abag(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get qpqt(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get wzwz(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get abab(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get qpqp(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get wzww(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get abaa(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get qpqq(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get wwxx(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get aarr(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get qqss(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get wwxy(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get aarg(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get qqst(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get wwxz(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get aarb(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get qqsp(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get wwxw(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get aara(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get qqsq(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get wwyx(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get aagr(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get qqts(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get wwyy(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get aagg(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get qqtt(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get wwyz(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get aagb(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get qqtp(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get wwyw(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get aaga(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get qqtq(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get wwzx(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get aabr(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get qqps(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get wwzy(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get aabg(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get qqpt(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get wwzz(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get aabb(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get qqpp(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get wwzw(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get aaba(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get qqpq(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get wwwx(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get aaar(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get qqqs(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get wwwy(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get aaag(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get qqqt(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get wwwz(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get aaab(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get qqqp(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get wwww(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get aaaa(): Node<NumOrBoolToVec4<TNumOrBool>>;
    get qqqq(): Node<NumOrBoolToVec4<TNumOrBool>>;
}

declare module "../core/Node.js" {
    interface NumOrBoolExtensions<TNumOrBool extends NumOrBoolType>
        extends
            Swizzle1From1<TNumOrBool>,
            Swizzle2From1<TNumOrBool>,
            Swizzle3From1<TNumOrBool>,
            Swizzle4From1<TNumOrBool>
    {
    }

    interface NumOrBoolVec2Extensions<TNumOrBool extends NumOrBoolType>
        extends
            Swizzle1From2<TNumOrBool>,
            Swizzle2From2<TNumOrBool>,
            Swizzle3From2<TNumOrBool>,
            Swizzle4From2<TNumOrBool>
    {
    }

    interface NumOrBoolVec3Extensions<TNumOrBool extends NumOrBoolType>
        extends
            Swizzle1From3<TNumOrBool>,
            Swizzle2From3<TNumOrBool>,
            Swizzle3From3<TNumOrBool>,
            Swizzle4From3<TNumOrBool>
    {
    }

    interface ColorExtensions
        extends Swizzle1From3<"float">, Swizzle2From3<"float">, Swizzle3From3<"float">, Swizzle4From3<"float">
    {
    }

    interface NumOrBoolVec4Extensions<TNumOrBool extends NumOrBoolType>
        extends
            Swizzle1From4<TNumOrBool>,
            Swizzle2From4<TNumOrBool>,
            Swizzle3From4<TNumOrBool>,
            Swizzle4From4<TNumOrBool>
    {
    }
}

interface SetSwizzle1<TNumOrBool extends NumOrBoolType> {
    setX(value: NumOrBool<TNumOrBool>): Node<TNumOrBool>;
    setR(value: NumOrBool<TNumOrBool>): Node<TNumOrBool>;
    setS(value: NumOrBool<TNumOrBool>): Node<TNumOrBool>;
}

interface SetSwizzle2<TNumOrBool extends NumOrBoolType> {
    setX(value: NumOrBool<TNumOrBool>): Node<NumOrBoolToVec2<TNumOrBool>>;
    setR(value: NumOrBool<TNumOrBool>): Node<NumOrBoolToVec2<TNumOrBool>>;
    setS(value: NumOrBool<TNumOrBool>): Node<NumOrBoolToVec2<TNumOrBool>>;
    setY(value: NumOrBool<TNumOrBool>): Node<NumOrBoolToVec2<TNumOrBool>>;
    setG(value: NumOrBool<TNumOrBool>): Node<NumOrBoolToVec2<TNumOrBool>>;
    setT(value: NumOrBool<TNumOrBool>): Node<NumOrBoolToVec2<TNumOrBool>>;
    setXY(value: Node<NumOrBoolToVec2<TNumOrBool>> | NumOrBool<TNumOrBool>): Node<NumOrBoolToVec2<TNumOrBool>>;
    setRG(value: Node<NumOrBoolToVec2<TNumOrBool>> | NumOrBool<TNumOrBool>): Node<NumOrBoolToVec2<TNumOrBool>>;
    setST(value: Node<NumOrBoolToVec2<TNumOrBool>> | NumOrBool<TNumOrBool>): Node<NumOrBoolToVec2<TNumOrBool>>;
}

interface SetSwizzle3<TNumOrBool extends NumOrBoolType> {
    setX(value: NumOrBool<TNumOrBool>): Node<NumOrBoolToVec3<TNumOrBool>>;
    setR(value: NumOrBool<TNumOrBool>): Node<NumOrBoolToVec3<TNumOrBool>>;
    setS(value: NumOrBool<TNumOrBool>): Node<NumOrBoolToVec3<TNumOrBool>>;
    setY(value: NumOrBool<TNumOrBool>): Node<NumOrBoolToVec3<TNumOrBool>>;
    setG(value: NumOrBool<TNumOrBool>): Node<NumOrBoolToVec3<TNumOrBool>>;
    setT(value: NumOrBool<TNumOrBool>): Node<NumOrBoolToVec3<TNumOrBool>>;
    setZ(value: NumOrBool<TNumOrBool>): Node<NumOrBoolToVec3<TNumOrBool>>;
    setB(value: NumOrBool<TNumOrBool>): Node<NumOrBoolToVec3<TNumOrBool>>;
    setP(value: NumOrBool<TNumOrBool>): Node<NumOrBoolToVec3<TNumOrBool>>;
    setXY(value: Node<NumOrBoolToVec2<TNumOrBool>> | NumOrBool<TNumOrBool>): Node<NumOrBoolToVec3<TNumOrBool>>;
    setRG(value: Node<NumOrBoolToVec2<TNumOrBool>> | NumOrBool<TNumOrBool>): Node<NumOrBoolToVec3<TNumOrBool>>;
    setST(value: Node<NumOrBoolToVec2<TNumOrBool>> | NumOrBool<TNumOrBool>): Node<NumOrBoolToVec3<TNumOrBool>>;
    setYZ(value: Node<NumOrBoolToVec2<TNumOrBool>> | NumOrBool<TNumOrBool>): Node<NumOrBoolToVec3<TNumOrBool>>;
    setGB(value: Node<NumOrBoolToVec2<TNumOrBool>> | NumOrBool<TNumOrBool>): Node<NumOrBoolToVec3<TNumOrBool>>;
    setTP(value: Node<NumOrBoolToVec2<TNumOrBool>> | NumOrBool<TNumOrBool>): Node<NumOrBoolToVec3<TNumOrBool>>;
    setXYZ(value: Node<NumOrBoolToVec3<TNumOrBool>> | NumOrBool<TNumOrBool>): Node<NumOrBoolToVec3<TNumOrBool>>;
    setRGB(value: Node<NumOrBoolToVec3<TNumOrBool>> | NumOrBool<TNumOrBool>): Node<NumOrBoolToVec3<TNumOrBool>>;
    setSTP(value: Node<NumOrBoolToVec3<TNumOrBool>> | NumOrBool<TNumOrBool>): Node<NumOrBoolToVec3<TNumOrBool>>;
}

interface SetSwizzle4<TNumOrBool extends NumOrBoolType> {
    setX(value: NumOrBool<TNumOrBool>): Node<NumOrBoolToVec4<TNumOrBool>>;
    setR(value: NumOrBool<TNumOrBool>): Node<NumOrBoolToVec4<TNumOrBool>>;
    setS(value: NumOrBool<TNumOrBool>): Node<NumOrBoolToVec4<TNumOrBool>>;
    setY(value: NumOrBool<TNumOrBool>): Node<NumOrBoolToVec4<TNumOrBool>>;
    setG(value: NumOrBool<TNumOrBool>): Node<NumOrBoolToVec4<TNumOrBool>>;
    setT(value: NumOrBool<TNumOrBool>): Node<NumOrBoolToVec4<TNumOrBool>>;
    setZ(value: NumOrBool<TNumOrBool>): Node<NumOrBoolToVec4<TNumOrBool>>;
    setB(value: NumOrBool<TNumOrBool>): Node<NumOrBoolToVec4<TNumOrBool>>;
    setP(value: NumOrBool<TNumOrBool>): Node<NumOrBoolToVec4<TNumOrBool>>;
    setW(value: NumOrBool<TNumOrBool>): Node<NumOrBoolToVec4<TNumOrBool>>;
    setA(value: NumOrBool<TNumOrBool>): Node<NumOrBoolToVec4<TNumOrBool>>;
    setQ(value: NumOrBool<TNumOrBool>): Node<NumOrBoolToVec4<TNumOrBool>>;
    setXY(value: Node<NumOrBoolToVec2<TNumOrBool>> | NumOrBool<TNumOrBool>): Node<NumOrBoolToVec4<TNumOrBool>>;
    setRG(value: Node<NumOrBoolToVec2<TNumOrBool>> | NumOrBool<TNumOrBool>): Node<NumOrBoolToVec4<TNumOrBool>>;
    setST(value: Node<NumOrBoolToVec2<TNumOrBool>> | NumOrBool<TNumOrBool>): Node<NumOrBoolToVec4<TNumOrBool>>;
    setYZ(value: Node<NumOrBoolToVec2<TNumOrBool>> | NumOrBool<TNumOrBool>): Node<NumOrBoolToVec4<TNumOrBool>>;
    setGB(value: Node<NumOrBoolToVec2<TNumOrBool>> | NumOrBool<TNumOrBool>): Node<NumOrBoolToVec4<TNumOrBool>>;
    setTP(value: Node<NumOrBoolToVec2<TNumOrBool>> | NumOrBool<TNumOrBool>): Node<NumOrBoolToVec4<TNumOrBool>>;
    setZW(value: Node<NumOrBoolToVec2<TNumOrBool>> | NumOrBool<TNumOrBool>): Node<NumOrBoolToVec4<TNumOrBool>>;
    setBA(value: Node<NumOrBoolToVec2<TNumOrBool>> | NumOrBool<TNumOrBool>): Node<NumOrBoolToVec4<TNumOrBool>>;
    setPQ(value: Node<NumOrBoolToVec2<TNumOrBool>> | NumOrBool<TNumOrBool>): Node<NumOrBoolToVec4<TNumOrBool>>;
    setXYZ(value: Node<NumOrBoolToVec3<TNumOrBool>> | NumOrBool<TNumOrBool>): Node<NumOrBoolToVec4<TNumOrBool>>;
    setRGB(value: Node<NumOrBoolToVec3<TNumOrBool>> | NumOrBool<TNumOrBool>): Node<NumOrBoolToVec4<TNumOrBool>>;
    setSTP(value: Node<NumOrBoolToVec3<TNumOrBool>> | NumOrBool<TNumOrBool>): Node<NumOrBoolToVec4<TNumOrBool>>;
    setYZW(value: Node<NumOrBoolToVec3<TNumOrBool>> | NumOrBool<TNumOrBool>): Node<NumOrBoolToVec4<TNumOrBool>>;
    setGBA(value: Node<NumOrBoolToVec3<TNumOrBool>> | NumOrBool<TNumOrBool>): Node<NumOrBoolToVec4<TNumOrBool>>;
    setTPQ(value: Node<NumOrBoolToVec3<TNumOrBool>> | NumOrBool<TNumOrBool>): Node<NumOrBoolToVec4<TNumOrBool>>;
    setXYZW(value: Node<NumOrBoolToVec4<TNumOrBool>> | NumOrBool<TNumOrBool>): Node<NumOrBoolToVec4<TNumOrBool>>;
    setRGBA(value: Node<NumOrBoolToVec4<TNumOrBool>> | NumOrBool<TNumOrBool>): Node<NumOrBoolToVec4<TNumOrBool>>;
    setSTPQ(value: Node<NumOrBoolToVec4<TNumOrBool>> | NumOrBool<TNumOrBool>): Node<NumOrBoolToVec4<TNumOrBool>>;
}

declare module "../core/Node.js" {
    // eslint-disable-next-line @typescript-eslint/no-empty-interface
    interface NumOrBoolExtensions<TNumOrBool extends NumOrBoolType> extends SetSwizzle1<TNumOrBool> {
    }

    // eslint-disable-next-line @typescript-eslint/no-empty-interface
    interface NumOrBoolVec2Extensions<TNumOrBool extends NumOrBoolType> extends SetSwizzle2<TNumOrBool> {
    }

    // eslint-disable-next-line @typescript-eslint/no-empty-interface
    interface NumOrBoolVec3Extensions<TNumOrBool extends NumOrBoolType> extends SetSwizzle3<TNumOrBool> {
    }

    // eslint-disable-next-line @typescript-eslint/no-empty-interface
    interface ColorExtensions extends SetSwizzle3<"float"> {
    }

    // eslint-disable-next-line @typescript-eslint/no-empty-interface
    interface NumOrBoolVec4Extensions<TNumOrBool extends NumOrBoolType> extends SetSwizzle4<TNumOrBool> {
    }
}

interface FlipSwizzle1 {
    flipX(): Node<"float">;
    flipR(): Node<"float">;
    flipS(): Node<"float">;
}

interface FlipSwizzle2 {
    flipX(): Node<"vec2">;
    flipR(): Node<"vec2">;
    flipS(): Node<"vec2">;
    flipY(): Node<"vec2">;
    flipG(): Node<"vec2">;
    flipT(): Node<"vec2">;
    flipXY(): Node<"vec2">;
    flipRG(): Node<"vec2">;
    flipST(): Node<"vec2">;
}

interface FlipSwizzle3 {
    flipX(): Node<"vec3">;
    flipR(): Node<"vec3">;
    flipS(): Node<"vec3">;
    flipY(): Node<"vec3">;
    flipG(): Node<"vec3">;
    flipT(): Node<"vec3">;
    flipZ(): Node<"vec3">;
    flipB(): Node<"vec3">;
    flipP(): Node<"vec3">;
    flipXY(): Node<"vec3">;
    flipRG(): Node<"vec3">;
    flipST(): Node<"vec3">;
    flipXZ(): Node<"vec3">;
    flipRB(): Node<"vec3">;
    flipSP(): Node<"vec3">;
    flipYZ(): Node<"vec3">;
    flipGB(): Node<"vec3">;
    flipTP(): Node<"vec3">;
    flipXYZ(): Node<"vec3">;
    flipRGB(): Node<"vec3">;
    flipSTP(): Node<"vec3">;
}

interface FlipSwizzle4 {
    flipX(): Node<"vec4">;
    flipR(): Node<"vec4">;
    flipS(): Node<"vec4">;
    flipY(): Node<"vec4">;
    flipG(): Node<"vec4">;
    flipT(): Node<"vec4">;
    flipZ(): Node<"vec4">;
    flipB(): Node<"vec4">;
    flipP(): Node<"vec4">;
    flipW(): Node<"vec4">;
    flipA(): Node<"vec4">;
    flipQ(): Node<"vec4">;
    flipXY(): Node<"vec4">;
    flipRG(): Node<"vec4">;
    flipST(): Node<"vec4">;
    flipXZ(): Node<"vec4">;
    flipRB(): Node<"vec4">;
    flipSP(): Node<"vec4">;
    flipYZ(): Node<"vec4">;
    flipGB(): Node<"vec4">;
    flipTP(): Node<"vec4">;
    flipXW(): Node<"vec4">;
    flipRA(): Node<"vec4">;
    flipSQ(): Node<"vec4">;
    flipYW(): Node<"vec4">;
    flipGA(): Node<"vec4">;
    flipTQ(): Node<"vec4">;
    flipZW(): Node<"vec4">;
    flipBA(): Node<"vec4">;
    flipPQ(): Node<"vec4">;
    flipXYZ(): Node<"vec4">;
    flipRGB(): Node<"vec4">;
    flipSTP(): Node<"vec4">;
    flipXYW(): Node<"vec4">;
    flipRGA(): Node<"vec4">;
    flipSTQ(): Node<"vec4">;
    flipXZW(): Node<"vec4">;
    flipRBA(): Node<"vec4">;
    flipSPQ(): Node<"vec4">;
    flipYZW(): Node<"vec4">;
    flipGBA(): Node<"vec4">;
    flipTPQ(): Node<"vec4">;
    flipXYZW(): Node<"vec4">;
    flipRGBA(): Node<"vec4">;
    flipSTPQ(): Node<"vec4">;
}

declare module "../core/Node.js" {
    // eslint-disable-next-line @typescript-eslint/no-empty-interface
    interface FloatExtensions extends FlipSwizzle1 {
    }

    // eslint-disable-next-line @typescript-eslint/no-empty-interface
    interface Vec2Extensions extends FlipSwizzle2 {
    }

    // eslint-disable-next-line @typescript-eslint/no-empty-interface
    interface Vec3Extensions extends FlipSwizzle3 {
    }

    // eslint-disable-next-line @typescript-eslint/no-empty-interface
    interface ColorExtensions extends FlipSwizzle3 {
    }

    // eslint-disable-next-line @typescript-eslint/no-empty-interface
    interface Vec4Extensions extends FlipSwizzle4 {
    }
}

export type NodeObject<T> = T extends Node ? T
    : T extends number ? Node<"float">
    : T extends boolean ? Node<"bool">
    : T extends (...args: never) => unknown ? unknown // FIXME This should return an FnNode
    : T extends Vector2 ? Node<"vec2">
    : T extends Vector3 ? Node<"vec3">
    : T extends Vector4 ? Node<"vec4">
    : T extends Matrix2 ? Node<"mat2">
    : T extends Matrix3 ? Node<"mat3">
    : T extends Matrix4 ? Node<"mat4">
    : T extends Color ? Node<"color">
    : T extends ArrayBuffer ? Node<"ArrayBuffer">
    : T;

type Proxied<T> = T extends Node<infer TNodeType>
    ? Node<TNodeType> extends T ? TNodeType extends "float" ? Node<"float"> | Node<"uint"> | number // FIXME remove Node<"uint">
        : TNodeType extends "bool" ? Node<"bool"> | boolean
        : TNodeType extends "vec2" ? Node<"vec2"> | Vector2
        : TNodeType extends "vec3" ? Node<"vec3"> | Vector3
        : TNodeType extends "vec4" ? Node<"vec4"> | Vector4
        : TNodeType extends "mat2" ? Node<"mat2"> | Matrix2
        : TNodeType extends "mat3" ? Node<"mat3"> | Matrix3
        : TNodeType extends "mat4" ? Node<"mat4"> | Matrix4
        : TNodeType extends "color" ? Node<"color"> | Color
        : TNodeType extends "ArrayBuffer" ? Node<"ArrayBuffer"> | ArrayBuffer
        : T
    : T
    : T;

type CoercibleToNodeType<TNodeType> = TNodeType extends "float" ? Node<"float"> | Node<"int"> | Node<"uint"> | number
    : TNodeType extends "bool" ? Node<"bool"> | boolean
    : TNodeType extends "vec2" ? Node<"vec2"> | Node<"ivec2"> | Node<"uvec2"> | Vector2
    : TNodeType extends "vec3" ? Node<"vec3"> | Node<"ivec3"> | Node<"uvec3"> | Vector3
    : TNodeType extends "vec4" ? Node<"vec4"> | Node<"ivec4"> | Node<"uvec4"> | Vector4
    : TNodeType extends "mat2" ? Node<"mat2"> | Matrix2
    : TNodeType extends "vec3" ? Node<"mat3"> | Matrix3
    : TNodeType extends "vec4" ? Node<"mat4"> | Matrix4
    : TNodeType extends "color" ? Node<"color"> | Color
    : TNodeType extends "ArrayBuffer" ? Node<"ArrayBuffer"> | ArrayBuffer
    : Node<TNodeType>;

// https://github.com/microsoft/TypeScript/issues/42435#issuecomment-765557874
// eslint-disable-next-line @definitelytyped/no-single-element-tuple-type
export type ProxiedTuple<T extends readonly [...unknown[]]> = [...{ [Index in keyof T]: Proxied<T[Index]> }];
export type ProxiedObject<T> = { [Index in keyof T]: Proxied<T[Index]> };
// eslint-disable-next-line @definitelytyped/no-single-element-tuple-type
type RemoveTail<T extends readonly [...unknown[]]> = T extends [unknown, ...infer X] ? X : [];
// eslint-disable-next-line @definitelytyped/no-single-element-tuple-type
type RemoveHeadAndTail<T extends readonly [...unknown[]]> = T extends [unknown, ...infer X, unknown] ? X : [];

/**
 * Temporary type to save signatures of 4 constructors. Each element may be tuple or undefined.
 *
 * We use an object instead of tuple or union as it makes stuff easier, especially in Typescript 4.0.
 */
interface Constructors<
    // eslint-disable-next-line @definitelytyped/no-single-element-tuple-type
    A extends undefined | [...unknown[]],
    // eslint-disable-next-line @definitelytyped/no-single-element-tuple-type
    B extends undefined | [...unknown[]],
    // eslint-disable-next-line @definitelytyped/no-single-element-tuple-type
    C extends undefined | [...unknown[]],
    // eslint-disable-next-line @definitelytyped/no-single-element-tuple-type
    D extends undefined | [...unknown[]],
> {
    a: A;
    b: B;
    c: C;
    d: D;
}

/**
 * Returns all constructors
 *
 * <https://github.com/microsoft/TypeScript/issues/37079>
 * <https://stackoverflow.com/a/52761156/1623826>
 */
type OverloadedConstructorsOf<T> = T extends {
    new(...args: infer A1): unknown;
    new(...args: infer A2): unknown;
    new(...args: infer A3): unknown;
    new(...args: infer A4): unknown;
} ? Constructors<A1, A2, A3, A4>
    : T extends {
        new(...args: infer A1): unknown;
        new(...args: infer A2): unknown;
        new(...args: infer A3): unknown;
    } ? Constructors<A1, A2, A3, undefined>
    : T extends {
        new(...args: infer A1): unknown;
        new(...args: infer A2): unknown;
    } ? Constructors<A1, A2, undefined, undefined>
    : T extends new(...args: infer A) => unknown ? Constructors<A, undefined, undefined, undefined>
    : Constructors<undefined, undefined, undefined, undefined>;

type AnyConstructors = Constructors<any, any, any, any>;

/**
 * Returns all constructors where the first parameter is assignable to given "scope"
 */
// eslint-disable-next-line @typescript-eslint/consistent-type-definitions
type FilterConstructorsByScope<T extends AnyConstructors, S> = {
    a: S extends T["a"][0] ? T["a"] : undefined;
    b: S extends T["b"][0] ? T["b"] : undefined;
    c: S extends T["c"][0] ? T["c"] : undefined;
    d: S extends T["d"][0] ? T["d"] : undefined;
};
/**
 * "flattens" the tuple into an union type
 */
type ConstructorUnion<T extends AnyConstructors> =
    | Exclude<T["a"], undefined>
    | Exclude<T["b"], undefined>
    | Exclude<T["c"], undefined>
    | Exclude<T["d"], undefined>;

/**
 * Extract list of possible scopes - union of the first parameter
 * of all constructors, should it be string
 */
type ExtractScopes<T extends AnyConstructors> =
    | (T["a"][0] extends string ? T["a"][0] : never)
    | (T["b"][0] extends string ? T["b"][0] : never)
    | (T["c"][0] extends string ? T["c"][0] : never)
    | (T["d"][0] extends string ? T["d"][0] : never);

type GetConstructorsByScope<T, S> = ConstructorUnion<FilterConstructorsByScope<OverloadedConstructorsOf<T>, S>>;
type GetConstructors<T> = ConstructorUnion<OverloadedConstructorsOf<T>>;
type GetPossibleScopes<T> = ExtractScopes<OverloadedConstructorsOf<T>>;

type NodeArray<T extends unknown[]> = { [Index in keyof T]: NodeObject<T[Index]> };
type NodeObjects<T> = { [Key in keyof T]: NodeObject<T[Key]> };
type ConstructedNode<T> = T extends new(...args: any[]) => infer R ? (R extends Node ? R : never) : never;

export type NodeOrType = Node | string;

type ShaderCallNodeInternal<TNodeType> = Node<TNodeType>;

type ShaderNodeInternal<TNodeType> = Node<TNodeType>;

export const defined: (v: unknown) => unknown;

export const getConstNodeType: (value: NodeOrType) => string | null;

export class ShaderNode<T = {}, R extends Node = Node> {
    constructor(jsFunc: (inputs: NodeObjects<T>, builder: NodeBuilder) => R);
    call: (
        inputs: { [key in keyof T]: T[key] extends Node ? Node : T[key] },
        builder?: NodeBuilder,
    ) => R;
}

export function nodeObject<T>(obj: T): NodeObject<T>;
export function nodeObjectIntent<T>(obj: T): NodeObject<T>;
export function nodeObjects<T>(obj: T): NodeObjects<T>;

// eslint-disable-next-line @definitelytyped/no-single-element-tuple-type
export function nodeArray<T extends unknown[]>(obj: readonly [...T]): NodeArray<T>;

export function nodeProxy<T>(
    nodeClass: T,
): (...params: ProxiedTuple<GetConstructors<T>>) => ConstructedNode<T>;

export function nodeProxy<T, S extends GetPossibleScopes<T>>(
    nodeClass: T,
    scope: S,
): (...params: ProxiedTuple<RemoveTail<GetConstructorsByScope<T, S>>>) => ConstructedNode<T>;

export function nodeProxy<T, S extends GetPossibleScopes<T>>(
    nodeClass: T,
    scope: S,
    factor: unknown,
): (...params: ProxiedTuple<RemoveHeadAndTail<GetConstructorsByScope<T, S>>>) => ConstructedNode<T>;

export function nodeImmutable<T>(
    nodeClass: T,
    ...params: ProxiedTuple<GetConstructors<T>>
): ConstructedNode<T>;

export function nodeProxyIntent<T>(
    nodeClass: T,
): (...params: ProxiedTuple<GetConstructors<T>>) => ConstructedNode<T>;

export function nodeProxyIntent<T, S extends GetPossibleScopes<T>>(
    nodeClass: T,
    scope: S,
): (...params: ProxiedTuple<RemoveTail<GetConstructorsByScope<T, S>>>) => ConstructedNode<T>;

export function nodeProxyIntent<T, S extends GetPossibleScopes<T>>(
    nodeClass: T,
    scope: S,
    factor: unknown,
): (...params: ProxiedTuple<RemoveHeadAndTail<GetConstructorsByScope<T, S>>>) => ConstructedNode<T>;

export const nodeProxyConstructor: unknown;

interface FullLayout {
    name: string;
    type: string;
    inputs: {
        name: string;
        type: string;
        qualifier?: "in" | "out" | "inout";
    }[];
}

interface AbbreviatedLayout {
    [inputName: string]: string;
    return: string;
}

type Layout = FullLayout | AbbreviatedLayout;

export interface FnNode<TArgs extends readonly unknown[], TReturn> {
    // eslint-disable-next-line @typescript-eslint/no-invalid-void-type
    (...args: TArgs): TReturn extends void ? ShaderCallNodeInternal<void> : TReturn;

    shaderNode: ShaderNodeInternal<TReturn>;
    id: number;

    setLayout(layout: Layout): this;

    generateNodeType(builder: NodeBuilder): string;

    once(subBuilds?: string[] | null): this;
}

export function Fn<TReturn>(
    jsFunc: (builder: NodeBuilder) => TReturn,
    layout?: Layout | string,
): FnNode<[], TReturn>;
export function Fn<TArgs extends readonly unknown[], TReturn>(
    jsFunc: (args: TArgs, builder: NodeBuilder) => TReturn,
    layout?: Layout | string,
): FnNode<ProxiedTuple<TArgs>, TReturn>;
export function Fn<TArgs extends { readonly [key: string]: unknown }, TReturn>(
    jsFunc: (args: TArgs, builder: NodeBuilder) => TReturn,
    layout?: Layout | string,
    // eslint-disable-next-line @definitelytyped/no-single-element-tuple-type
): FnNode<[ProxiedObject<TArgs>], TReturn>;

export const setCurrentStack: (stack: StackNode | null) => void;

export const getCurrentStack: () => StackNode | null;

export const If: (boolNode: Node, method: () => void) => StackNode;
export const Switch: (expression: Node) => StackNode;

export function Stack(node: Node): Node;

declare module "../core/Node.js" {
    interface NodeElements {
        toStack: () => Node;
    }
}

/**
 * Can be implicitly converted to a float
 */
type ScalarNode = Node<"float"> | Node<"int"> | Node<"uint"> | Node<"bool">;

/**
 * Can be implicitly converted to a float
 */
type Scalar = ScalarNode | number | boolean;

interface ColorFunction {
    // The first branch in `ConvertType` will forward the parameters to the `Color` constructor if there are no
    //   parameters or all the parameters are non-objects
    (color?: string | number): VarNode<"color", ConstNode<"color", Color>>;
    (r: number, g: number, b: number): VarNode<"color", ConstNode<"color", Color>>;

    // The second branch does not apply because `cacheMap` is `null`

    // The third branch will be triggered if there is a single parameter
    // ConstNode
    (color: Color | Vector3): VarNode<"color", ConstNode<"color", Color>>;
    // ConvertNode
    (node: Node<"float">): VarNode<"color", ConvertNode<"color">>;
    (node: Node<"color"> | Node<"vec3">): VarNode<"color", ConvertNode<"color">>;

    // The fall-through branch will be triggered if there is more than one parameter, and one of the parameters is an
    //   object
    (
        x: Node<"float"> | number,
        y: Node<"float"> | number,
        z: Node<"float"> | number,
    ): VarNode<"vec3", JoinNode<"vec3">>;
    (xy: Node<"vec2"> | Vector2, z: Node<"float"> | number): VarNode<"vec3", JoinNode<"vec3">>;
    (x: Node<"float"> | number, yz: Node<"vec2"> | Vector2): VarNode<"vec3", JoinNode<"vec3">>;
}

export const color: ColorFunction;

interface FloatFunction {
    // ConstNode
    (value?: number): VarNode<"float", ConstNode<"float", number>>;

    // ConvertNode
    (node: ScalarNode): VarNode<"float", ConvertNode<"float">>;
}

export const float: FloatFunction;

interface IntFunction {
    // ConstNode
    (value?: number): VarNode<"int", ConstNode<"int", number>>;

    // ConvertNode
    (node: ScalarNode): VarNode<"int", ConvertNode<"int">>;
}

export const int: IntFunction;

interface UintFunction {
    // ConstNode
    (value?: number): VarNode<"uint", ConstNode<"uint", number>>;

    // ConvertNode
    (node: ScalarNode): VarNode<"uint", ConvertNode<"uint">>;
}

export const uint: UintFunction;

interface BoolFunction {
    // ConstNode
    (value?: boolean): VarNode<"bool", ConstNode<"bool", boolean>>;

    // ConvertNode
    (node: ScalarNode): VarNode<"bool", ConvertNode<"bool">>;
}

export const bool: BoolFunction;

interface Vec2Function {
    // The first branch in `ConvertType` will forward the parameters to the `Vector2` constructor if there are no
    //   parameters or all the parameters are non-objects
    (x?: number, y?: number): VarNode<"vec2", ConstNode<"vec2", Vector2>>;

    // The second branch does not apply because `cacheMap` is `null`

    // The third branch will be triggered if there is a single parameter
    // ConstNode
    (value: Vector2): VarNode<"vec2", ConstNode<"vec2", Vector2>>;
    // ConvertNode
    (node: ScalarNode): VarNode<"vec2", ConvertNode<"vec2">>;
    (node: Node<"vec2"> | Node<"ivec2"> | Node<"uvec2"> | Node<"bvec2">): VarNode<"vec2", ConvertNode<"vec2">>;

    // The fall-through branch will be triggered if there is more than one parameter, and one of the parameters is an
    //   object
    (x: Scalar, y: Scalar): VarNode<"vec2", JoinNode<"vec2">>;
}

export const vec2: Vec2Function;

interface Ivec2Function {
    // The first branch in `ConvertType` will forward the parameters to the `Vector2` constructor if there are no
    //   parameters or all the parameters are non-objects
    (x?: number, y?: number): VarNode<"ivec2", ConstNode<"ivec2", Vector2>>;

    // The second branch does not apply because `cacheMap` is `null`

    // The third branch will be triggered if there is a single parameter
    // ConstNode
    (value: Vector2): VarNode<"ivec2", ConstNode<"ivec2", Vector2>>;
    // ConvertNode
    (node: Node<"int">): VarNode<"ivec2", ConvertNode<"ivec2">>;
    (node: Node<"vec2"> | Node<"ivec2"> | Node<"uvec2"> | Node<"bvec2">): VarNode<"ivec2", ConvertNode<"ivec2">>;

    // The fall-through branch will be triggered if there is more than one parameter, and one of the parameters is an
    //   object
    (x: Node<"int"> | number, y: Node<"int"> | number): VarNode<"ivec2", JoinNode<"ivec2">>;
}

export const ivec2: Ivec2Function;

interface Uvec2Function {
    // The first branch in `ConvertType` will forward the parameters to the `Vector2` constructor if there are no
    //   parameters or all the parameters are non-objects
    (x?: number, y?: number): VarNode<"uvec2", ConstNode<"uvec2", Vector2>>;

    // The second branch does not apply because `cacheMap` is `null`

    // The third branch will be triggered if there is a single parameter
    // ConstNode
    (value: Vector2): VarNode<"uvec2", ConstNode<"uvec2", Vector2>>;
    // ConvertNode
    (node: Node<"uint">): VarNode<"uvec2", ConvertNode<"uvec2">>;
    (node: Node<"vec2"> | Node<"ivec2"> | Node<"uvec2"> | Node<"bvec2">): VarNode<"uvec2", ConvertNode<"uvec2">>;

    // The fall-through branch will be triggered if there is more than one parameter, and one of the parameters is an
    //   object
    (x: Node<"uint"> | number, y: Node<"uint"> | number): VarNode<"uvec2", JoinNode<"uvec2">>;
}

export const uvec2: Uvec2Function;

interface Bvec2Function {
    // The first branch in `ConvertType` will forward the parameters to the `Vector2` constructor if there are no
    //   parameters or all the parameters are non-objects
    (x?: boolean, y?: boolean): VarNode<"bvec2", ConstNode<"bvec2", Vector2>>;

    // The second branch does not apply because `cacheMap` is `null`

    // The third branch will be triggered if there is a single parameter
    // ConstNode
    (value: Vector2): VarNode<"bvec2", ConstNode<"bvec2", Vector2>>;
    // ConvertNode
    (node: Node<"bool">): VarNode<"bvec2", ConvertNode<"bvec2">>;
    (node: Node<"vec2"> | Node<"ivec2"> | Node<"uvec2"> | Node<"bvec2">): VarNode<"bvec2", ConvertNode<"bvec2">>;

    // The fall-through branch will be triggered if there is more than one parameter, and one of the parameters is an
    //   object
    (x: Node<"bool"> | boolean, y: Node<"bool"> | boolean): VarNode<"bvec2", JoinNode<"bvec2">>;
}

export const bvec2: Bvec2Function;

interface Vec3Function {
    // The first branch in `ConvertType` will forward the parameters to the `Vector3` constructor if there are no
    //   parameters or all the parameters are non-objects
    (x?: number, y?: number, z?: number): VarNode<"vec3", ConstNode<"vec3", Vector3>>;

    // The second branch does not apply because `cacheMap` is `null`

    // The third branch will be triggered if there is a single parameter
    // ConstNode
    (value: Vector3): VarNode<"vec3", ConstNode<"vec3", Vector3>>;
    // ConvertNode
    (node: ScalarNode): VarNode<"vec3", ConvertNode<"vec3">>;
    (
        node: Node<"vec3"> | Node<"ivec3"> | Node<"uvec3"> | Node<"bvec3"> | Node<"vec4">,
    ): VarNode<"vec3", ConvertNode<"vec3">>;

    // The fall-through branch will be triggered if there is more than one parameter, and one of the parameters is an
    //   object
    (x: Scalar, y: Scalar, z: Scalar): VarNode<"vec3", JoinNode<"vec3">>;
    (xy: Node<"vec2"> | Vector2, z: Scalar): VarNode<"vec3", JoinNode<"vec3">>;
    (x: Scalar, yz: Node<"vec2"> | Vector2): VarNode<"vec3", JoinNode<"vec3">>;
}

export const vec3: Vec3Function;

interface Ivec3Function {
    // The first branch in `ConvertType` will forward the parameters to the `Vector3` constructor if there are no
    //   parameters or all the parameters are non-objects
    (x?: number, y?: number, z?: number): VarNode<"ivec3", ConstNode<"ivec3", Vector3>>;

    // The second branch does not apply because `cacheMap` is `null`

    // The third branch will be triggered if there is a single parameter
    // ConstNode
    (value: Vector3): VarNode<"ivec3", ConstNode<"ivec3", Vector3>>;
    // ConvertNode
    (node: Node<"int">): VarNode<"ivec3", ConvertNode<"ivec3">>;
    (node: Node<"vec3"> | Node<"ivec3"> | Node<"uvec3"> | Node<"bvec3">): VarNode<"ivec3", ConvertNode<"ivec3">>;

    // The fall-through branch will be triggered if there is more than one parameter, and one of the parameters is an
    //   object
    (x: Node<"int"> | number, y: Node<"int"> | number, z: Node<"int"> | number): VarNode<"ivec3", JoinNode<"ivec3">>;
    (xy: Node<"ivec2"> | Vector2, z: Node<"int"> | number): VarNode<"ivec3", JoinNode<"ivec3">>;
    (x: Node<"int"> | number, yz: Node<"ivec2"> | Vector2): VarNode<"ivec3", JoinNode<"ivec3">>;
}

export const ivec3: Ivec3Function;

interface Uvec3Function {
    // The first branch in `ConvertType` will forward the parameters to the `Vector3` constructor if there are no
    //   parameters or all the parameters are non-objects
    (x?: number, y?: number, z?: number): VarNode<"uvec3", ConstNode<"uvec3", Vector3>>;

    // The second branch does not apply because `cacheMap` is `null`

    // The third branch will be triggered if there is a single parameter
    // ConstNode
    (value: Vector3): VarNode<"uvec3", ConstNode<"uvec3", Vector3>>;
    // ConvertNode
    (node: Node<"uint">): VarNode<"uvec3", ConvertNode<"uvec3">>;
    (node: Node<"vec3"> | Node<"ivec3"> | Node<"uvec3"> | Node<"bvec3">): VarNode<"uvec3", ConvertNode<"uvec3">>;

    // The fall-through branch will be triggered if there is more than one parameter, and one of the parameters is an
    //   object
    (x: Node<"uint"> | number, y: Node<"uint"> | number, z: Node<"uint"> | number): VarNode<"uvec3", JoinNode<"uvec3">>;
    (xy: Node<"uvec2"> | Vector2, z: Node<"uint"> | number): VarNode<"uvec3", JoinNode<"uvec3">>;
    (x: Node<"uint"> | number, yz: Node<"uvec2"> | Vector2): VarNode<"uvec3", JoinNode<"uvec3">>;
}

export const uvec3: Uvec3Function;

interface Bvec3Function {
    // The first branch in `ConvertType` will forward the parameters to the `Vector3` constructor if there are no
    //   parameters or all the parameters are non-objects
    (x?: boolean, y?: boolean, z?: boolean): VarNode<"bvec3", ConstNode<"bvec3", Vector3>>;

    // The second branch does not apply because `cacheMap` is `null`

    // The third branch will be triggered if there is a single parameter
    // ConstNode
    (value: Vector3): VarNode<"bvec3", ConstNode<"bvec3", Vector3>>;
    // ConvertNode
    (node: Node<"bool">): VarNode<"bvec3", ConvertNode<"bvec3">>;
    (node: Node<"vec3"> | Node<"ivec3"> | Node<"uvec3"> | Node<"bvec3">): VarNode<"bvec3", ConvertNode<"bvec3">>;

    // The fall-through branch will be triggered if there is more than one parameter, and one of the parameters is an
    //   object
    (
        x: Node<"bool"> | boolean,
        y: Node<"bool"> | boolean,
        z: Node<"bool"> | boolean,
    ): VarNode<"bvec3", JoinNode<"bvec3">>;
    (xy: Node<"bvec2"> | Vector2, z: Node<"bool"> | boolean): VarNode<"bvec3", JoinNode<"bvec3">>;
    (x: Node<"bool"> | boolean, yz: Node<"bvec2"> | Vector2): VarNode<"bvec3", JoinNode<"bvec3">>;
}

export const bvec3: Bvec3Function;

interface Vec4Function {
    // The first branch in `ConvertType` will forward the parameters to the `Vector4` constructor if there are no
    //   parameters or all the parameters are non-objects
    (x?: number, y?: number, z?: number, w?: number): VarNode<"vec4", ConstNode<"vec4", Vector4>>;

    // The second branch does not apply because `cacheMap` is `null`

    // The third branch will be triggered if there is a single parameter
    // ConstNode
    (value: Vector4): VarNode<"vec4", ConstNode<"vec4", Vector4>>;
    // ConvertNode
    (node: ScalarNode): VarNode<"vec4", ConvertNode<"vec4">>;
    (node: Node<"vec4"> | Node<"ivec4"> | Node<"uvec4"> | Node<"bvec4">): VarNode<"vec4", ConvertNode<"vec4">>;

    // The fall-through branch will be triggered if there is more than one parameter, and one of the parameters is an
    //   object
    (x: Scalar, y: Scalar, z: Scalar, w: Scalar): VarNode<"vec4", JoinNode<"vec4">>;
    (x: Scalar, yz: Node<"vec2"> | Vector2, w: Scalar): VarNode<"vec4", JoinNode<"vec4">>;
    (x: Scalar, y: Scalar, zw: Node<"vec2"> | Vector2): VarNode<"vec4", JoinNode<"vec4">>;
    (xy: Node<"vec2"> | Vector2, zw: Node<"vec2"> | Vector2): VarNode<"vec4", JoinNode<"vec4">>;
    (xy: Node<"vec2"> | Vector2, z: Scalar, w: Scalar): VarNode<"vec4", JoinNode<"vec4">>;
    (xyz: Node<"vec3"> | Node<"color"> | Vector3 | Color | Node<"vec4">, w: Scalar): VarNode<"vec4", JoinNode<"vec4">>;
    (x: Scalar, yzw: Node<"vec3"> | Node<"color"> | Vector3 | Color): VarNode<"vec4", JoinNode<"vec4">>;
}

export const vec4: Vec4Function;

interface Ivec4Function {
    // The first branch in `ConvertType` will forward the parameters to the `Vector4` constructor if there are no
    //   parameters or all the parameters are non-objects
    (x?: number, y?: number, z?: number, w?: number): VarNode<"ivec4", ConstNode<"ivec4", Vector4>>;

    // The second branch does not apply because `cacheMap` is `null`

    // The third branch will be triggered if there is a single parameter
    // ConstNode
    (value: Vector4): VarNode<"ivec4", ConstNode<"ivec4", Vector4>>;
    // ConvertNode
    (node: Node<"int">): VarNode<"ivec4", ConvertNode<"ivec4">>;
    (node: Node<"vec4"> | Node<"ivec4"> | Node<"uvec4"> | Node<"bvec4">): VarNode<"ivec4", ConvertNode<"ivec4">>;

    // The fall-through branch will be triggered if there is more than one parameter, and one of the parameters is an
    //   object
    (
        x: Node<"int"> | number,
        y: Node<"int"> | number,
        z: Node<"int"> | number,
        w: Node<"int"> | number,
    ): VarNode<"ivec4", JoinNode<"ivec4">>;
    (
        x: Node<"int"> | number,
        yz: Node<"ivec2"> | Vector2,
        w: Node<"int"> | number,
    ): VarNode<"ivec4", JoinNode<"ivec4">>;
    (
        x: Node<"int"> | number,
        y: Node<"int"> | number,
        zw: Node<"ivec2"> | Vector2,
    ): VarNode<"ivec4", JoinNode<"ivec4">>;
    (xy: Node<"ivec2"> | Vector2, zw: Node<"ivec2"> | Vector2): VarNode<"ivec4", JoinNode<"ivec4">>;
    (
        xy: Node<"ivec2"> | Vector2,
        z: Node<"int"> | number,
        w: Node<"int"> | number,
    ): VarNode<"ivec4", JoinNode<"ivec4">>;
    (xyz: Node<"ivec3"> | Vector3, w: Node<"int"> | number): VarNode<"ivec4", JoinNode<"ivec4">>;
    (x: Node<"int"> | number, yzw: Node<"ivec3"> | Vector3): VarNode<"ivec4", JoinNode<"ivec4">>;
}

export const ivec4: Ivec4Function;

interface Uvec4Function {
    // The first branch in `ConvertType` will forward the parameters to the `Vector4` constructor if there are no
    //   parameters or all the parameters are non-objects
    (x?: number, y?: number, z?: number, w?: number): VarNode<"uvec4", ConstNode<"uvec4", Vector4>>;

    // The second branch does not apply because `cacheMap` is `null`

    // The third branch will be triggered if there is a single parameter
    // ConstNode
    (value: Vector4): VarNode<"uvec4", ConstNode<"uvec4", Vector4>>;
    // ConvertNode
    (node: Node<"uint">): VarNode<"uvec4", ConvertNode<"uvec4">>;
    (node: Node<"vec4"> | Node<"ivec4"> | Node<"uvec4"> | Node<"bvec4">): VarNode<"uvec4", ConvertNode<"uvec4">>;

    // The fall-through branch will be triggered if there is more than one parameter, and one of the parameters is an
    //   object
    (
        x: Node<"uint"> | number,
        y: Node<"uint"> | number,
        z: Node<"uint"> | number,
        w: Node<"uint"> | number,
    ): VarNode<"uvec4", JoinNode<"uvec4">>;
    (
        x: Node<"uint"> | number,
        yz: Node<"uvec2"> | Vector2,
        w: Node<"uint"> | number,
    ): VarNode<"uvec4", JoinNode<"uvec4">>;
    (
        x: Node<"uint"> | number,
        y: Node<"uint"> | number,
        zw: Node<"uvec2"> | Vector2,
    ): VarNode<"uvec4", JoinNode<"uvec4">>;
    (xy: Node<"uvec2"> | Vector2, zw: Node<"uvec2"> | Vector2): VarNode<"uvec4", JoinNode<"uvec4">>;
    (
        xy: Node<"uvec2"> | Vector2,
        z: Node<"uint"> | number,
        w: Node<"uint"> | number,
    ): VarNode<"uvec4", JoinNode<"uvec4">>;
    (xyz: Node<"uvec3"> | Vector3, w: Node<"uint"> | number): VarNode<"uvec4", JoinNode<"uvec4">>;
    (x: Node<"uint"> | number, yzw: Node<"uvec3"> | Vector3): VarNode<"uvec4", JoinNode<"uvec4">>;
}

export const uvec4: Uvec4Function;

interface Bvec4Function {
    // The first branch in `ConvertType` will forward the parameters to the `Vector4` constructor if there are no
    //   parameters or all the parameters are non-objects
    (x?: boolean, y?: boolean, z?: boolean, w?: boolean): VarNode<"bvec4", ConstNode<"bvec4", Vector4>>;

    // The second branch does not apply because `cacheMap` is `null`

    // The third branch will be triggered if there is a single parameter
    // ConstNode
    (value: Vector4): VarNode<"bvec4", ConstNode<"bvec4", Vector4>>;
    // ConvertNode
    (node: Node<"bool">): VarNode<"bvec4", ConvertNode<"bvec4">>;
    (node: Node<"vec4"> | Node<"ivec4"> | Node<"uvec4"> | Node<"bvec4">): VarNode<"bvec4", ConvertNode<"bvec4">>;

    // The fall-through branch will be triggered if there is more than one parameter, and one of the parameters is an
    //   object
    (
        x: Node<"bool"> | boolean,
        y: Node<"bool"> | boolean,
        z: Node<"bool"> | boolean,
        w: Node<"bool"> | boolean,
    ): VarNode<"bvec4", JoinNode<"bvec4">>;
    (
        x: Node<"bool"> | boolean,
        yz: Node<"bvec2"> | Vector2,
        w: Node<"bool"> | boolean,
    ): VarNode<"bvec4", JoinNode<"bvec4">>;
    (
        x: Node<"bool"> | boolean,
        y: Node<"bool"> | boolean,
        zw: Node<"bvec2"> | Vector2,
    ): VarNode<"bvec4", JoinNode<"bvec4">>;
    (xy: Node<"bvec2"> | Vector2, zw: Node<"bvec2"> | Vector2): VarNode<"bvec4", JoinNode<"bvec4">>;
    (
        xy: Node<"bvec2"> | Vector2,
        z: Node<"bool"> | boolean,
        w: Node<"bool"> | boolean,
    ): VarNode<"bvec4", JoinNode<"bvec4">>;
    (xyz: Node<"bvec3"> | Vector3, w: Node<"bool"> | boolean): VarNode<"bvec4", JoinNode<"bvec4">>;
    (x: Node<"bool"> | boolean, yzw: Node<"bvec3"> | Vector3): VarNode<"bvec4", JoinNode<"bvec4">>;
}

export const bvec4: Bvec4Function;

interface Mat2Function {
    // The first branch in `ConvertType` will forward the parameters to the `Matrix2` constructor if there are no
    //   parameters or all the parameters are non-objects
    (): VarNode<"mat2", ConstNode<"mat2", Matrix2>>;
    (n11: number, n12: number, n21: number, n22: number): VarNode<"mat2", ConstNode<"mat2", Matrix2>>;

    // The second branch does not apply because `cacheMap` is `null`

    // The third branch will be triggered if there is a single parameter
    // ConstNode
    (value: Matrix2): VarNode<"mat2", ConstNode<"mat2", Matrix2>>;
    // ConvertNode
    (node: Node<"mat2">): VarNode<"mat2", ConvertNode<"mat2">>;

    // The fall-through branch will be triggered if there is more than one parameter, and one of the parameters is an
    //   object
    (n1: Node<"vec2"> | Vector2, n2: Node<"vec2"> | Vector2): VarNode<"mat2", JoinNode<"mat2">>;
    (
        n11: Node<"float"> | number,
        n12: Node<"float"> | number,
        n21: Node<"float"> | number,
        n22: Node<"float"> | number,
    ): VarNode<"mat2", JoinNode<"mat2">>;
}

export const mat2: Mat2Function;

interface Mat3Function {
    // The first branch in `ConvertType` will forward the parameters to the `Matrix3` constructor if there are no
    //   parameters or all the parameters are non-objects
    (): VarNode<"mat3", ConstNode<"mat3", Matrix3>>;
    (
        n11: number,
        n12: number,
        n13: number,
        n21: number,
        n22: number,
        n23: number,
        n31: number,
        n32: number,
        n33: number,
    ): VarNode<"mat3", ConstNode<"mat3", Matrix3>>;

    // The second branch does not apply because `cacheMap` is `null`

    // The third branch will be triggered if there is a single parameter
    // ConstNode
    (value: Matrix3): VarNode<"mat3", ConstNode<"mat3", Matrix3>>;
    // ConvertNode
    (node: Node<"mat3">): VarNode<"mat3", ConvertNode<"mat3">>;

    // The fall-through branch will be triggered if there is more than one parameter, and one of the parameters is an
    //   object
    (
        n1: Node<"vec3"> | Vector3,
        n2: Node<"vec3"> | Vector3,
        n3: Node<"vec3"> | Vector3,
    ): VarNode<"mat3", JoinNode<"mat3">>;
    (
        n11: Node<"float"> | number,
        n12: Node<"float"> | number,
        n13: Node<"float"> | number,
        n21: Node<"float"> | number,
        n22: Node<"float"> | number,
        n23: Node<"float"> | number,
        n31: Node<"float"> | number,
        n32: Node<"float"> | number,
        n33: Node<"float"> | number,
    ): VarNode<"mat3", JoinNode<"mat3">>;
}

export const mat3: Mat3Function;

interface Mat4Function {
    // The first branch in `ConvertType` will forward the parameters to the `Matrix4` constructor if there are no
    //   parameters or all the parameters are non-objects
    (): VarNode<"mat4", ConstNode<"mat4", Matrix4>>;
    (
        n11: number,
        n12: number,
        n13: number,
        n14: number,
        n21: number,
        n22: number,
        n23: number,
        n24: number,
        n31: number,
        n32: number,
        n33: number,
        n34: number,
        n41: number,
        n42: number,
        n43: number,
        n44: number,
    ): VarNode<"mat4", ConstNode<"mat4", Matrix4>>;

    // The second branch does not apply because `cacheMap` is `null`

    // The third branch will be triggered if there is a single parameter
    // ConstNode
    (value: Matrix4): VarNode<"mat4", ConstNode<"mat4", Matrix4>>;
    // ConvertNode
    (node: Node<"mat4">): VarNode<"mat4", ConvertNode<"mat4">>;

    // The fall-through branch will be triggered if there is more than one parameter, and one of the parameters is an
    //   object
    (
        n1: Node<"vec4"> | Vector4,
        n2: Node<"vec4"> | Vector4,
        n3: Node<"vec4"> | Vector4,
        n4: Node<"vec4"> | Vector4,
    ): VarNode<"mat4", JoinNode<"mat4">>;
    (
        n11: Node<"float"> | number,
        n12: Node<"float"> | number,
        n13: Node<"float"> | number,
        n14: Node<"float"> | number,
        n21: Node<"float"> | number,
        n22: Node<"float"> | number,
        n23: Node<"float"> | number,
        n24: Node<"float"> | number,
        n31: Node<"float"> | number,
        n32: Node<"float"> | number,
        n33: Node<"float"> | number,
        n34: Node<"float"> | number,
        n41: Node<"float"> | number,
        n42: Node<"float"> | number,
        n43: Node<"float"> | number,
        n44: Node<"float"> | number,
    ): VarNode<"mat4", JoinNode<"mat4">>;
}

export const mat4: Mat4Function;

declare module "../core/Node.js" {
    interface ColorExtensions {
        toColor: () => VarNode<"color", ConvertNode<"color">>;
    }

    interface NumOrBoolExtensions<TNumOrBool extends NumOrBoolType> {
        toFloat: () => VarNode<"float", ConvertNode<"float">>;
        toInt: () => VarNode<"int", ConvertNode<"int">>;
        toUint: () => VarNode<"uint", ConvertNode<"uint">>;
        toBool: () => VarNode<"bool", ConvertNode<"bool">>;
    }

    interface NumOrBoolVec2Extensions<TNumOrBool extends NumOrBoolType> {
        toVec2: () => VarNode<"vec2", ConvertNode<"vec2">>;
        toIVec2: () => VarNode<"ivec2", ConvertNode<"ivec2">>;
        toUVec2: () => VarNode<"uvec2", ConvertNode<"uvec2">>;
        toBVec2: () => VarNode<"bvec2", ConvertNode<"bvec2">>;
    }

    interface NumOrBoolVec3Extensions<TNumOrBool extends NumOrBoolType> {
        toColor: () => VarNode<"color", ConvertNode<"color">>;
        toVec3: () => VarNode<"vec3", ConvertNode<"vec3">>;
        toIVec3: () => VarNode<"ivec3", ConvertNode<"ivec3">>;
        toUVec3: () => VarNode<"uvec3", ConvertNode<"uvec3">>;
        toBVec3: () => VarNode<"bvec3", ConvertNode<"bvec3">>;
    }

    interface NumOrBoolVec4Extensions<TNumOrBool extends NumOrBoolType> {
        toVec4: () => VarNode<"vec4", ConvertNode<"vec4">>;
        toIVec4: () => VarNode<"ivec4", ConvertNode<"ivec4">>;
        toUVec4: () => VarNode<"uvec4", ConvertNode<"uvec4">>;
        toBVec4: () => VarNode<"bvec4", ConvertNode<"bvec4">>;
    }

    interface MatExtensions<TMat extends MatType> {
        toMat2: () => VarNode<"mat2", ConvertNode<"mat2">>;
        toMat3: () => VarNode<"mat3", ConvertNode<"mat3">>;
        toMat4: () => VarNode<"mat4", ConvertNode<"mat4">>;
    }
}

export const element: <TNodeType>(node: ArrayNode<TNodeType>, indexNode: Node | number) => ArrayElementNode<TNodeType>;
export const convert: (node: Node, types: string) => Node;
export const split: (node: Node, channels?: string) => Node;

declare module "../core/ArrayNode.js" {
    interface ArrayNodeInterface<TNodeType> {
        element: (indexNode: Node | number) => ArrayElementNode<TNodeType>;
    }
}

declare module "../core/Node.js" {
    interface NodeElements {
        convert: (types: string) => Node;
    }
}

/**
 * @deprecated append() has been renamed to Stack().
 */
export const append: (node: Node) => Node;

declare module "../core/Node.js" {
    interface NodeElements {
        /**
         * @deprecated append() has been renamed to Stack().
         */
        append: () => Node;
    }
}

export {};

