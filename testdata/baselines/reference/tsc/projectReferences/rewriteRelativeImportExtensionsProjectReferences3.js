
currentDirectory::/home/src/workspaces/solution
useCaseSensitiveFileNames::true
Input::--p src/services --pretty false
//// [/home/src/workspaces/solution/dist/compiler/parser.d.ts] new file
export {};
//// [/home/src/workspaces/solution/src/compiler/parser.ts] new file
export {};
//// [/home/src/workspaces/solution/src/compiler/tsconfig.json] new file
{
					"extends": "../tsconfig-base.json",
					"compilerOptions": {
						"rootDir": ".",
						"outDir": "../../dist/compiler"
					}
				}
//// [/home/src/workspaces/solution/src/services/services.ts] new file
import {} from "../compiler/parser.ts";
//// [/home/src/workspaces/solution/src/services/tsconfig.json] new file
{
					"extends": "../tsconfig-base.json",
					"compilerOptions": {
						"rootDir": ".", 
						"outDir": "../../dist/services"
					},
					"references": [
						{ "path": "../compiler" }
					]
				}
//// [/home/src/workspaces/solution/src/tsconfig-base.json] new file
{
					"compilerOptions": { 
						"module": "nodenext",
						"composite": true,
						"rewriteRelativeImportExtensions": true
					}
				}

ExitStatus:: 0

CompilerOptions::{
    "project": "/home/src/workspaces/solution/src/services",
    "pretty": false
}
Output::
No output
//// [/home/src/workspaces/solution/dist/compiler/parser.d.ts] no change
//// [/home/src/workspaces/solution/dist/services/services.d.ts] new file
export {};

//// [/home/src/workspaces/solution/dist/services/services.js] new file
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });

//// [/home/src/workspaces/solution/dist/services/tsconfig.tsbuildinfo] new file
{"version":"FakeTSVersion","fileNames":["bundled:///libs/lib.es5.d.ts","bundled:///libs/lib.es2015.d.ts","bundled:///libs/lib.es2016.d.ts","bundled:///libs/lib.es2017.d.ts","bundled:///libs/lib.es2018.d.ts","bundled:///libs/lib.es2019.d.ts","bundled:///libs/lib.es2020.d.ts","bundled:///libs/lib.es2021.d.ts","bundled:///libs/lib.es2022.d.ts","bundled:///libs/lib.es2023.d.ts","bundled:///libs/lib.es2024.d.ts","bundled:///libs/lib.esnext.d.ts","bundled:///libs/lib.dom.d.ts","bundled:///libs/lib.dom.iterable.d.ts","bundled:///libs/lib.dom.asynciterable.d.ts","bundled:///libs/lib.webworker.importscripts.d.ts","bundled:///libs/lib.scripthost.d.ts","bundled:///libs/lib.es2015.core.d.ts","bundled:///libs/lib.es2015.collection.d.ts","bundled:///libs/lib.es2015.generator.d.ts","bundled:///libs/lib.es2015.iterable.d.ts","bundled:///libs/lib.es2015.promise.d.ts","bundled:///libs/lib.es2015.proxy.d.ts","bundled:///libs/lib.es2015.reflect.d.ts","bundled:///libs/lib.es2015.symbol.d.ts","bundled:///libs/lib.es2015.symbol.wellknown.d.ts","bundled:///libs/lib.es2016.array.include.d.ts","bundled:///libs/lib.es2016.intl.d.ts","bundled:///libs/lib.es2017.arraybuffer.d.ts","bundled:///libs/lib.es2017.date.d.ts","bundled:///libs/lib.es2017.object.d.ts","bundled:///libs/lib.es2017.sharedmemory.d.ts","bundled:///libs/lib.es2017.string.d.ts","bundled:///libs/lib.es2017.intl.d.ts","bundled:///libs/lib.es2017.typedarrays.d.ts","bundled:///libs/lib.es2018.asyncgenerator.d.ts","bundled:///libs/lib.es2018.asynciterable.d.ts","bundled:///libs/lib.es2018.intl.d.ts","bundled:///libs/lib.es2018.promise.d.ts","bundled:///libs/lib.es2018.regexp.d.ts","bundled:///libs/lib.es2019.array.d.ts","bundled:///libs/lib.es2019.object.d.ts","bundled:///libs/lib.es2019.string.d.ts","bundled:///libs/lib.es2019.symbol.d.ts","bundled:///libs/lib.es2019.intl.d.ts","bundled:///libs/lib.es2020.bigint.d.ts","bundled:///libs/lib.es2020.date.d.ts","bundled:///libs/lib.es2020.promise.d.ts","bundled:///libs/lib.es2020.sharedmemory.d.ts","bundled:///libs/lib.es2020.string.d.ts","bundled:///libs/lib.es2020.symbol.wellknown.d.ts","bundled:///libs/lib.es2020.intl.d.ts","bundled:///libs/lib.es2020.number.d.ts","bundled:///libs/lib.es2021.promise.d.ts","bundled:///libs/lib.es2021.string.d.ts","bundled:///libs/lib.es2021.weakref.d.ts","bundled:///libs/lib.es2021.intl.d.ts","bundled:///libs/lib.es2022.array.d.ts","bundled:///libs/lib.es2022.error.d.ts","bundled:///libs/lib.es2022.intl.d.ts","bundled:///libs/lib.es2022.object.d.ts","bundled:///libs/lib.es2022.string.d.ts","bundled:///libs/lib.es2022.regexp.d.ts","bundled:///libs/lib.es2023.array.d.ts","bundled:///libs/lib.es2023.collection.d.ts","bundled:///libs/lib.es2023.intl.d.ts","bundled:///libs/lib.es2024.arraybuffer.d.ts","bundled:///libs/lib.es2024.collection.d.ts","bundled:///libs/lib.es2024.object.d.ts","bundled:///libs/lib.es2024.promise.d.ts","bundled:///libs/lib.es2024.regexp.d.ts","bundled:///libs/lib.es2024.sharedmemory.d.ts","bundled:///libs/lib.es2024.string.d.ts","bundled:///libs/lib.esnext.array.d.ts","bundled:///libs/lib.esnext.collection.d.ts","bundled:///libs/lib.esnext.intl.d.ts","bundled:///libs/lib.esnext.disposable.d.ts","bundled:///libs/lib.esnext.promise.d.ts","bundled:///libs/lib.esnext.decorators.d.ts","bundled:///libs/lib.esnext.iterator.d.ts","bundled:///libs/lib.esnext.float16.d.ts","bundled:///libs/lib.decorators.d.ts","bundled:///libs/lib.decorators.legacy.d.ts","bundled:///libs/lib.esnext.full.d.ts","../compiler/parser.d.ts","../../src/services/services.ts"],"fileInfos":[{"version":"69684132aeb9b5642cbcd9e22dff7818ff0ee1aa831728af0ecf97d3364d5546","affectsGlobalScope":true,"impliedNodeFormat":1},"45b7ab580deca34ae9729e97c13cfd999df04416a79116c3bfb483804f85ded4","3facaf05f0c5fc569c5649dd359892c98a85557e3e0c847964caeb67076f4d75","e44bb8bbac7f10ecc786703fe0a6a4b952189f908707980ba8f3c8975a760962","5e1c4c362065a6b95ff952c0eab010f04dcd2c3494e813b493ecfd4fcb9fc0d8","68d73b4a11549f9c0b7d352d10e91e5dca8faa3322bfb77b661839c42b1ddec7","5efce4fc3c29ea84e8928f97adec086e3dc876365e0982cc8479a07954a3efd4","feecb1be483ed332fad555aff858affd90a48ab19ba7272ee084704eb7167569","ee7bad0c15b58988daa84371e0b89d313b762ab83cb5b31b8a2d1162e8eb41c2","27bdc30a0e32783366a5abeda841bc22757c1797de8681bbe81fbc735eeb1c10","8fd575e12870e9944c7e1d62e1f5a73fcf23dd8d3a321f2a2c74c20d022283fe","8bf8b5e44e3c9c36f98e1007e8b7018c0f38d8adc07aecef42f5200114547c70",{"version":"092c2bfe125ce69dbb1223c85d68d4d2397d7d8411867b5cc03cec902c233763","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"07f073f19d67f74d732b1adea08e1dc66b1b58d77cb5b43931dee3d798a2fd53","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"d7a3c8b952931daebdfc7a2897c53c0a1c73624593fa070e46bd537e64dcd20a","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"80e18897e5884b6723488d4f5652167e7bb5024f946743134ecc4aa4ee731f89","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"cd034f499c6cdca722b60c04b5b1b78e058487a7085a8e0d6fb50809947ee573","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"c57796738e7f83dbc4b8e65132f11a377649c00dd3eee333f672b8f0a6bea671","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"dc2df20b1bcdc8c2d34af4926e2c3ab15ffe1160a63e58b7e09833f616efff44","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"515d0b7b9bea2e31ea4ec968e9edd2c39d3eebf4a2d5cbd04e88639819ae3b71","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"0559b1f683ac7505ae451f9a96ce4c3c92bdc71411651ca6ddb0e88baaaad6a3","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"0dc1e7ceda9b8b9b455c3a2d67b0412feab00bd2f66656cd8850e8831b08b537","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"ce691fb9e5c64efb9547083e4a34091bcbe5bdb41027e310ebba8f7d96a98671","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"8d697a2a929a5fcb38b7a65594020fcef05ec1630804a33748829c5ff53640d0","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"4ff2a353abf8a80ee399af572debb8faab2d33ad38c4b4474cff7f26e7653b8d","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"936e80ad36a2ee83fc3caf008e7c4c5afe45b3cf3d5c24408f039c1d47bdc1df","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"d15bea3d62cbbdb9797079416b8ac375ae99162a7fba5de2c6c505446486ac0a","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"68d18b664c9d32a7336a70235958b8997ebc1c3b8505f4f1ae2b7e7753b87618","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"eb3d66c8327153d8fa7dd03f9c58d351107fe824c79e9b56b462935176cdf12a","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"38f0219c9e23c915ef9790ab1d680440d95419ad264816fa15009a8851e79119","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"69ab18c3b76cd9b1be3d188eaf8bba06112ebbe2f47f6c322b5105a6fbc45a2e","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"fef8cfad2e2dc5f5b3d97a6f4f2e92848eb1b88e897bb7318cef0e2820bceaab","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"2f11ff796926e0832f9ae148008138ad583bd181899ab7dd768a2666700b1893","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"4de680d5bb41c17f7f68e0419412ca23c98d5749dcaaea1896172f06435891fc","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"954296b30da6d508a104a3a0b5d96b76495c709785c1d11610908e63481ee667","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"ac9538681b19688c8eae65811b329d3744af679e0bdfa5d842d0e32524c73e1c","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"0a969edff4bd52585473d24995c5ef223f6652d6ef46193309b3921d65dd4376","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"9e9fbd7030c440b33d021da145d3232984c8bb7916f277e8ffd3dc2e3eae2bdb","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"811ec78f7fefcabbda4bfa93b3eb67d9ae166ef95f9bff989d964061cbf81a0c","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"717937616a17072082152a2ef351cb51f98802fb4b2fdabd32399843875974ca","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"d7e7d9b7b50e5f22c915b525acc5a49a7a6584cf8f62d0569e557c5cfc4b2ac2","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"71c37f4c9543f31dfced6c7840e068c5a5aacb7b89111a4364b1d5276b852557","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"576711e016cf4f1804676043e6a0a5414252560eb57de9faceee34d79798c850","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"89c1b1281ba7b8a96efc676b11b264de7a8374c5ea1e6617f11880a13fc56dc6","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"74f7fa2d027d5b33eb0471c8e82a6c87216223181ec31247c357a3e8e2fddc5b","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"d6d7ae4d1f1f3772e2a3cde568ed08991a8ae34a080ff1151af28b7f798e22ca","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"063600664504610fe3e99b717a1223f8b1900087fab0b4cad1496a114744f8df","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"934019d7e3c81950f9a8426d093458b65d5aff2c7c1511233c0fd5b941e608ab","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"52ada8e0b6e0482b728070b7639ee42e83a9b1c22d205992756fe020fd9f4a47","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"3bdefe1bfd4d6dee0e26f928f93ccc128f1b64d5d501ff4a8cf3c6371200e5e6","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"59fb2c069260b4ba00b5643b907ef5d5341b167e7d1dbf58dfd895658bda2867","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"639e512c0dfc3fad96a84caad71b8834d66329a1f28dc95e3946c9b58176c73a","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"368af93f74c9c932edd84c58883e736c9e3d53cec1fe24c0b0ff451f529ceab1","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"af3dd424cf267428f30ccfc376f47a2c0114546b55c44d8c0f1d57d841e28d74","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"995c005ab91a498455ea8dfb63aa9f83fa2ea793c3d8aa344be4a1678d06d399","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"959d36cddf5e7d572a65045b876f2956c973a586da58e5d26cde519184fd9b8a","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"965f36eae237dd74e6cca203a43e9ca801ce38824ead814728a2807b1910117d","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"3925a6c820dcb1a06506c90b1577db1fdbf7705d65b62b99dce4be75c637e26b","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"0a3d63ef2b853447ec4f749d3f368ce642264246e02911fcb1590d8c161b8005","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"b5ce7a470bc3628408429040c4e3a53a27755022a32fd05e2cb694e7015386c7","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"8444af78980e3b20b49324f4a16ba35024fef3ee069a0eb67616ea6ca821c47a","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"3287d9d085fbd618c3971944b65b4be57859f5415f495b33a6adc994edd2f004","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"b4b67b1a91182421f5df999988c690f14d813b9850b40acd06ed44691f6727ad","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"df83c2a6c73228b625b0beb6669c7ee2a09c914637e2d35170723ad49c0f5cd4","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"436aaf437562f276ec2ddbee2f2cdedac7664c1e4c1d2c36839ddd582eeb3d0a","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"8e3c06ea092138bf9fa5e874a1fdbc9d54805d074bee1de31b99a11e2fec239d","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"87dc0f382502f5bbce5129bdc0aea21e19a3abbc19259e0b43ae038a9fc4e326","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"b1cb28af0c891c8c96b2d6b7be76bd394fddcfdb4709a20ba05a7c1605eea0f9","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"2fef54945a13095fdb9b84f705f2b5994597640c46afeb2ce78352fab4cb3279","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"ac77cb3e8c6d3565793eb90a8373ee8033146315a3dbead3bde8db5eaf5e5ec6","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"56e4ed5aab5f5920980066a9409bfaf53e6d21d3f8d020c17e4de584d29600ad","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"4ece9f17b3866cc077099c73f4983bddbcb1dc7ddb943227f1ec070f529dedd1","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"0a6282c8827e4b9a95f4bf4f5c205673ada31b982f50572d27103df8ceb8013c","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"1c9319a09485199c1f7b0498f2988d6d2249793ef67edda49d1e584746be9032","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"e3a2a0cee0f03ffdde24d89660eba2685bfbdeae955a6c67e8c4c9fd28928eeb","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"811c71eee4aa0ac5f7adf713323a5c41b0cf6c4e17367a34fbce379e12bbf0a4","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"51ad4c928303041605b4d7ae32e0c1ee387d43a24cd6f1ebf4a2699e1076d4fa","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"60037901da1a425516449b9a20073aa03386cce92f7a1fd902d7602be3a7c2e9","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"d4b1d2c51d058fc21ec2629fff7a76249dec2e36e12960ea056e3ef89174080f","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"22adec94ef7047a6c9d1af3cb96be87a335908bf9ef386ae9fd50eeb37f44c47","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"4245fee526a7d1754529d19227ecbf3be066ff79ebb6a380d78e41648f2f224d","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"8e7f8264d0fb4c5339605a15daadb037bf238c10b654bb3eee14208f860a32ea","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"782dec38049b92d4e85c1585fbea5474a219c6984a35b004963b00beb1aab538","affectsGlobalScope":true,"impliedNodeFormat":1},"bde31fd423cd93b0eff97197a3f66df7c93e8c0c335cbeb113b7ff1ac35c23f4","2e29cd9a98755c46896f7a2d56524db2d6d96b248e36db46de14c30bf47c8d05",{"version":"407537635fda1a543a422ecdd456c1402aaa2083cde5acfb4eb424ab02fc0612","signature":"8e609bb71c20b858c77f0e9f90bb1319db8477b13f9f965f1a1e18524bf50881","impliedNodeFormat":1}],"fileIdsList":[[85]],"options":{"composite":true,"module":199,"outDir":"./","rewriteRelativeImportExtensions":true,"rootDir":"../../src/services"},"referencedMap":[[86,1]],"latestChangedDtsFile":"./services.d.ts"}
//// [/home/src/workspaces/solution/dist/services/tsconfig.tsbuildinfo.readable.baseline.txt] new file
{
  "version": "FakeTSVersion",
  "fileNames": [
    "bundled:///libs/lib.es5.d.ts",
    "bundled:///libs/lib.es2015.d.ts",
    "bundled:///libs/lib.es2016.d.ts",
    "bundled:///libs/lib.es2017.d.ts",
    "bundled:///libs/lib.es2018.d.ts",
    "bundled:///libs/lib.es2019.d.ts",
    "bundled:///libs/lib.es2020.d.ts",
    "bundled:///libs/lib.es2021.d.ts",
    "bundled:///libs/lib.es2022.d.ts",
    "bundled:///libs/lib.es2023.d.ts",
    "bundled:///libs/lib.es2024.d.ts",
    "bundled:///libs/lib.esnext.d.ts",
    "bundled:///libs/lib.dom.d.ts",
    "bundled:///libs/lib.dom.iterable.d.ts",
    "bundled:///libs/lib.dom.asynciterable.d.ts",
    "bundled:///libs/lib.webworker.importscripts.d.ts",
    "bundled:///libs/lib.scripthost.d.ts",
    "bundled:///libs/lib.es2015.core.d.ts",
    "bundled:///libs/lib.es2015.collection.d.ts",
    "bundled:///libs/lib.es2015.generator.d.ts",
    "bundled:///libs/lib.es2015.iterable.d.ts",
    "bundled:///libs/lib.es2015.promise.d.ts",
    "bundled:///libs/lib.es2015.proxy.d.ts",
    "bundled:///libs/lib.es2015.reflect.d.ts",
    "bundled:///libs/lib.es2015.symbol.d.ts",
    "bundled:///libs/lib.es2015.symbol.wellknown.d.ts",
    "bundled:///libs/lib.es2016.array.include.d.ts",
    "bundled:///libs/lib.es2016.intl.d.ts",
    "bundled:///libs/lib.es2017.arraybuffer.d.ts",
    "bundled:///libs/lib.es2017.date.d.ts",
    "bundled:///libs/lib.es2017.object.d.ts",
    "bundled:///libs/lib.es2017.sharedmemory.d.ts",
    "bundled:///libs/lib.es2017.string.d.ts",
    "bundled:///libs/lib.es2017.intl.d.ts",
    "bundled:///libs/lib.es2017.typedarrays.d.ts",
    "bundled:///libs/lib.es2018.asyncgenerator.d.ts",
    "bundled:///libs/lib.es2018.asynciterable.d.ts",
    "bundled:///libs/lib.es2018.intl.d.ts",
    "bundled:///libs/lib.es2018.promise.d.ts",
    "bundled:///libs/lib.es2018.regexp.d.ts",
    "bundled:///libs/lib.es2019.array.d.ts",
    "bundled:///libs/lib.es2019.object.d.ts",
    "bundled:///libs/lib.es2019.string.d.ts",
    "bundled:///libs/lib.es2019.symbol.d.ts",
    "bundled:///libs/lib.es2019.intl.d.ts",
    "bundled:///libs/lib.es2020.bigint.d.ts",
    "bundled:///libs/lib.es2020.date.d.ts",
    "bundled:///libs/lib.es2020.promise.d.ts",
    "bundled:///libs/lib.es2020.sharedmemory.d.ts",
    "bundled:///libs/lib.es2020.string.d.ts",
    "bundled:///libs/lib.es2020.symbol.wellknown.d.ts",
    "bundled:///libs/lib.es2020.intl.d.ts",
    "bundled:///libs/lib.es2020.number.d.ts",
    "bundled:///libs/lib.es2021.promise.d.ts",
    "bundled:///libs/lib.es2021.string.d.ts",
    "bundled:///libs/lib.es2021.weakref.d.ts",
    "bundled:///libs/lib.es2021.intl.d.ts",
    "bundled:///libs/lib.es2022.array.d.ts",
    "bundled:///libs/lib.es2022.error.d.ts",
    "bundled:///libs/lib.es2022.intl.d.ts",
    "bundled:///libs/lib.es2022.object.d.ts",
    "bundled:///libs/lib.es2022.string.d.ts",
    "bundled:///libs/lib.es2022.regexp.d.ts",
    "bundled:///libs/lib.es2023.array.d.ts",
    "bundled:///libs/lib.es2023.collection.d.ts",
    "bundled:///libs/lib.es2023.intl.d.ts",
    "bundled:///libs/lib.es2024.arraybuffer.d.ts",
    "bundled:///libs/lib.es2024.collection.d.ts",
    "bundled:///libs/lib.es2024.object.d.ts",
    "bundled:///libs/lib.es2024.promise.d.ts",
    "bundled:///libs/lib.es2024.regexp.d.ts",
    "bundled:///libs/lib.es2024.sharedmemory.d.ts",
    "bundled:///libs/lib.es2024.string.d.ts",
    "bundled:///libs/lib.esnext.array.d.ts",
    "bundled:///libs/lib.esnext.collection.d.ts",
    "bundled:///libs/lib.esnext.intl.d.ts",
    "bundled:///libs/lib.esnext.disposable.d.ts",
    "bundled:///libs/lib.esnext.promise.d.ts",
    "bundled:///libs/lib.esnext.decorators.d.ts",
    "bundled:///libs/lib.esnext.iterator.d.ts",
    "bundled:///libs/lib.esnext.float16.d.ts",
    "bundled:///libs/lib.decorators.d.ts",
    "bundled:///libs/lib.decorators.legacy.d.ts",
    "bundled:///libs/lib.esnext.full.d.ts",
    "../compiler/parser.d.ts",
    "../../src/services/services.ts"
  ],
  "fileInfos": [
    {
      "fileName": "bundled:///libs/lib.es5.d.ts",
      "version": "69684132aeb9b5642cbcd9e22dff7818ff0ee1aa831728af0ecf97d3364d5546",
      "signature": "69684132aeb9b5642cbcd9e22dff7818ff0ee1aa831728af0ecf97d3364d5546",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "69684132aeb9b5642cbcd9e22dff7818ff0ee1aa831728af0ecf97d3364d5546",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "bundled:///libs/lib.es2015.d.ts",
      "version": "45b7ab580deca34ae9729e97c13cfd999df04416a79116c3bfb483804f85ded4",
      "signature": "45b7ab580deca34ae9729e97c13cfd999df04416a79116c3bfb483804f85ded4",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "bundled:///libs/lib.es2016.d.ts",
      "version": "3facaf05f0c5fc569c5649dd359892c98a85557e3e0c847964caeb67076f4d75",
      "signature": "3facaf05f0c5fc569c5649dd359892c98a85557e3e0c847964caeb67076f4d75",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "bundled:///libs/lib.es2017.d.ts",
      "version": "e44bb8bbac7f10ecc786703fe0a6a4b952189f908707980ba8f3c8975a760962",
      "signature": "e44bb8bbac7f10ecc786703fe0a6a4b952189f908707980ba8f3c8975a760962",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "bundled:///libs/lib.es2018.d.ts",
      "version": "5e1c4c362065a6b95ff952c0eab010f04dcd2c3494e813b493ecfd4fcb9fc0d8",
      "signature": "5e1c4c362065a6b95ff952c0eab010f04dcd2c3494e813b493ecfd4fcb9fc0d8",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "bundled:///libs/lib.es2019.d.ts",
      "version": "68d73b4a11549f9c0b7d352d10e91e5dca8faa3322bfb77b661839c42b1ddec7",
      "signature": "68d73b4a11549f9c0b7d352d10e91e5dca8faa3322bfb77b661839c42b1ddec7",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "bundled:///libs/lib.es2020.d.ts",
      "version": "5efce4fc3c29ea84e8928f97adec086e3dc876365e0982cc8479a07954a3efd4",
      "signature": "5efce4fc3c29ea84e8928f97adec086e3dc876365e0982cc8479a07954a3efd4",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "bundled:///libs/lib.es2021.d.ts",
      "version": "feecb1be483ed332fad555aff858affd90a48ab19ba7272ee084704eb7167569",
      "signature": "feecb1be483ed332fad555aff858affd90a48ab19ba7272ee084704eb7167569",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "bundled:///libs/lib.es2022.d.ts",
      "version": "ee7bad0c15b58988daa84371e0b89d313b762ab83cb5b31b8a2d1162e8eb41c2",
      "signature": "ee7bad0c15b58988daa84371e0b89d313b762ab83cb5b31b8a2d1162e8eb41c2",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "bundled:///libs/lib.es2023.d.ts",
      "version": "27bdc30a0e32783366a5abeda841bc22757c1797de8681bbe81fbc735eeb1c10",
      "signature": "27bdc30a0e32783366a5abeda841bc22757c1797de8681bbe81fbc735eeb1c10",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "bundled:///libs/lib.es2024.d.ts",
      "version": "8fd575e12870e9944c7e1d62e1f5a73fcf23dd8d3a321f2a2c74c20d022283fe",
      "signature": "8fd575e12870e9944c7e1d62e1f5a73fcf23dd8d3a321f2a2c74c20d022283fe",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "bundled:///libs/lib.esnext.d.ts",
      "version": "8bf8b5e44e3c9c36f98e1007e8b7018c0f38d8adc07aecef42f5200114547c70",
      "signature": "8bf8b5e44e3c9c36f98e1007e8b7018c0f38d8adc07aecef42f5200114547c70",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "bundled:///libs/lib.dom.d.ts",
      "version": "092c2bfe125ce69dbb1223c85d68d4d2397d7d8411867b5cc03cec902c233763",
      "signature": "092c2bfe125ce69dbb1223c85d68d4d2397d7d8411867b5cc03cec902c233763",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "092c2bfe125ce69dbb1223c85d68d4d2397d7d8411867b5cc03cec902c233763",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "bundled:///libs/lib.dom.iterable.d.ts",
      "version": "07f073f19d67f74d732b1adea08e1dc66b1b58d77cb5b43931dee3d798a2fd53",
      "signature": "07f073f19d67f74d732b1adea08e1dc66b1b58d77cb5b43931dee3d798a2fd53",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "07f073f19d67f74d732b1adea08e1dc66b1b58d77cb5b43931dee3d798a2fd53",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "bundled:///libs/lib.dom.asynciterable.d.ts",
      "version": "d7a3c8b952931daebdfc7a2897c53c0a1c73624593fa070e46bd537e64dcd20a",
      "signature": "d7a3c8b952931daebdfc7a2897c53c0a1c73624593fa070e46bd537e64dcd20a",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "d7a3c8b952931daebdfc7a2897c53c0a1c73624593fa070e46bd537e64dcd20a",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "bundled:///libs/lib.webworker.importscripts.d.ts",
      "version": "80e18897e5884b6723488d4f5652167e7bb5024f946743134ecc4aa4ee731f89",
      "signature": "80e18897e5884b6723488d4f5652167e7bb5024f946743134ecc4aa4ee731f89",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "80e18897e5884b6723488d4f5652167e7bb5024f946743134ecc4aa4ee731f89",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "bundled:///libs/lib.scripthost.d.ts",
      "version": "cd034f499c6cdca722b60c04b5b1b78e058487a7085a8e0d6fb50809947ee573",
      "signature": "cd034f499c6cdca722b60c04b5b1b78e058487a7085a8e0d6fb50809947ee573",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "cd034f499c6cdca722b60c04b5b1b78e058487a7085a8e0d6fb50809947ee573",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "bundled:///libs/lib.es2015.core.d.ts",
      "version": "c57796738e7f83dbc4b8e65132f11a377649c00dd3eee333f672b8f0a6bea671",
      "signature": "c57796738e7f83dbc4b8e65132f11a377649c00dd3eee333f672b8f0a6bea671",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "c57796738e7f83dbc4b8e65132f11a377649c00dd3eee333f672b8f0a6bea671",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "bundled:///libs/lib.es2015.collection.d.ts",
      "version": "dc2df20b1bcdc8c2d34af4926e2c3ab15ffe1160a63e58b7e09833f616efff44",
      "signature": "dc2df20b1bcdc8c2d34af4926e2c3ab15ffe1160a63e58b7e09833f616efff44",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "dc2df20b1bcdc8c2d34af4926e2c3ab15ffe1160a63e58b7e09833f616efff44",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "bundled:///libs/lib.es2015.generator.d.ts",
      "version": "515d0b7b9bea2e31ea4ec968e9edd2c39d3eebf4a2d5cbd04e88639819ae3b71",
      "signature": "515d0b7b9bea2e31ea4ec968e9edd2c39d3eebf4a2d5cbd04e88639819ae3b71",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "515d0b7b9bea2e31ea4ec968e9edd2c39d3eebf4a2d5cbd04e88639819ae3b71",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "bundled:///libs/lib.es2015.iterable.d.ts",
      "version": "0559b1f683ac7505ae451f9a96ce4c3c92bdc71411651ca6ddb0e88baaaad6a3",
      "signature": "0559b1f683ac7505ae451f9a96ce4c3c92bdc71411651ca6ddb0e88baaaad6a3",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "0559b1f683ac7505ae451f9a96ce4c3c92bdc71411651ca6ddb0e88baaaad6a3",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "bundled:///libs/lib.es2015.promise.d.ts",
      "version": "0dc1e7ceda9b8b9b455c3a2d67b0412feab00bd2f66656cd8850e8831b08b537",
      "signature": "0dc1e7ceda9b8b9b455c3a2d67b0412feab00bd2f66656cd8850e8831b08b537",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "0dc1e7ceda9b8b9b455c3a2d67b0412feab00bd2f66656cd8850e8831b08b537",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "bundled:///libs/lib.es2015.proxy.d.ts",
      "version": "ce691fb9e5c64efb9547083e4a34091bcbe5bdb41027e310ebba8f7d96a98671",
      "signature": "ce691fb9e5c64efb9547083e4a34091bcbe5bdb41027e310ebba8f7d96a98671",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "ce691fb9e5c64efb9547083e4a34091bcbe5bdb41027e310ebba8f7d96a98671",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "bundled:///libs/lib.es2015.reflect.d.ts",
      "version": "8d697a2a929a5fcb38b7a65594020fcef05ec1630804a33748829c5ff53640d0",
      "signature": "8d697a2a929a5fcb38b7a65594020fcef05ec1630804a33748829c5ff53640d0",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "8d697a2a929a5fcb38b7a65594020fcef05ec1630804a33748829c5ff53640d0",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "bundled:///libs/lib.es2015.symbol.d.ts",
      "version": "4ff2a353abf8a80ee399af572debb8faab2d33ad38c4b4474cff7f26e7653b8d",
      "signature": "4ff2a353abf8a80ee399af572debb8faab2d33ad38c4b4474cff7f26e7653b8d",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "4ff2a353abf8a80ee399af572debb8faab2d33ad38c4b4474cff7f26e7653b8d",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "bundled:///libs/lib.es2015.symbol.wellknown.d.ts",
      "version": "936e80ad36a2ee83fc3caf008e7c4c5afe45b3cf3d5c24408f039c1d47bdc1df",
      "signature": "936e80ad36a2ee83fc3caf008e7c4c5afe45b3cf3d5c24408f039c1d47bdc1df",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "936e80ad36a2ee83fc3caf008e7c4c5afe45b3cf3d5c24408f039c1d47bdc1df",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "bundled:///libs/lib.es2016.array.include.d.ts",
      "version": "d15bea3d62cbbdb9797079416b8ac375ae99162a7fba5de2c6c505446486ac0a",
      "signature": "d15bea3d62cbbdb9797079416b8ac375ae99162a7fba5de2c6c505446486ac0a",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "d15bea3d62cbbdb9797079416b8ac375ae99162a7fba5de2c6c505446486ac0a",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "bundled:///libs/lib.es2016.intl.d.ts",
      "version": "68d18b664c9d32a7336a70235958b8997ebc1c3b8505f4f1ae2b7e7753b87618",
      "signature": "68d18b664c9d32a7336a70235958b8997ebc1c3b8505f4f1ae2b7e7753b87618",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "68d18b664c9d32a7336a70235958b8997ebc1c3b8505f4f1ae2b7e7753b87618",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "bundled:///libs/lib.es2017.arraybuffer.d.ts",
      "version": "eb3d66c8327153d8fa7dd03f9c58d351107fe824c79e9b56b462935176cdf12a",
      "signature": "eb3d66c8327153d8fa7dd03f9c58d351107fe824c79e9b56b462935176cdf12a",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "eb3d66c8327153d8fa7dd03f9c58d351107fe824c79e9b56b462935176cdf12a",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "bundled:///libs/lib.es2017.date.d.ts",
      "version": "38f0219c9e23c915ef9790ab1d680440d95419ad264816fa15009a8851e79119",
      "signature": "38f0219c9e23c915ef9790ab1d680440d95419ad264816fa15009a8851e79119",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "38f0219c9e23c915ef9790ab1d680440d95419ad264816fa15009a8851e79119",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "bundled:///libs/lib.es2017.object.d.ts",
      "version": "69ab18c3b76cd9b1be3d188eaf8bba06112ebbe2f47f6c322b5105a6fbc45a2e",
      "signature": "69ab18c3b76cd9b1be3d188eaf8bba06112ebbe2f47f6c322b5105a6fbc45a2e",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "69ab18c3b76cd9b1be3d188eaf8bba06112ebbe2f47f6c322b5105a6fbc45a2e",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "bundled:///libs/lib.es2017.sharedmemory.d.ts",
      "version": "fef8cfad2e2dc5f5b3d97a6f4f2e92848eb1b88e897bb7318cef0e2820bceaab",
      "signature": "fef8cfad2e2dc5f5b3d97a6f4f2e92848eb1b88e897bb7318cef0e2820bceaab",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "fef8cfad2e2dc5f5b3d97a6f4f2e92848eb1b88e897bb7318cef0e2820bceaab",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "bundled:///libs/lib.es2017.string.d.ts",
      "version": "2f11ff796926e0832f9ae148008138ad583bd181899ab7dd768a2666700b1893",
      "signature": "2f11ff796926e0832f9ae148008138ad583bd181899ab7dd768a2666700b1893",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "2f11ff796926e0832f9ae148008138ad583bd181899ab7dd768a2666700b1893",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "bundled:///libs/lib.es2017.intl.d.ts",
      "version": "4de680d5bb41c17f7f68e0419412ca23c98d5749dcaaea1896172f06435891fc",
      "signature": "4de680d5bb41c17f7f68e0419412ca23c98d5749dcaaea1896172f06435891fc",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "4de680d5bb41c17f7f68e0419412ca23c98d5749dcaaea1896172f06435891fc",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "bundled:///libs/lib.es2017.typedarrays.d.ts",
      "version": "954296b30da6d508a104a3a0b5d96b76495c709785c1d11610908e63481ee667",
      "signature": "954296b30da6d508a104a3a0b5d96b76495c709785c1d11610908e63481ee667",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "954296b30da6d508a104a3a0b5d96b76495c709785c1d11610908e63481ee667",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "bundled:///libs/lib.es2018.asyncgenerator.d.ts",
      "version": "ac9538681b19688c8eae65811b329d3744af679e0bdfa5d842d0e32524c73e1c",
      "signature": "ac9538681b19688c8eae65811b329d3744af679e0bdfa5d842d0e32524c73e1c",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "ac9538681b19688c8eae65811b329d3744af679e0bdfa5d842d0e32524c73e1c",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "bundled:///libs/lib.es2018.asynciterable.d.ts",
      "version": "0a969edff4bd52585473d24995c5ef223f6652d6ef46193309b3921d65dd4376",
      "signature": "0a969edff4bd52585473d24995c5ef223f6652d6ef46193309b3921d65dd4376",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "0a969edff4bd52585473d24995c5ef223f6652d6ef46193309b3921d65dd4376",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "bundled:///libs/lib.es2018.intl.d.ts",
      "version": "9e9fbd7030c440b33d021da145d3232984c8bb7916f277e8ffd3dc2e3eae2bdb",
      "signature": "9e9fbd7030c440b33d021da145d3232984c8bb7916f277e8ffd3dc2e3eae2bdb",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "9e9fbd7030c440b33d021da145d3232984c8bb7916f277e8ffd3dc2e3eae2bdb",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "bundled:///libs/lib.es2018.promise.d.ts",
      "version": "811ec78f7fefcabbda4bfa93b3eb67d9ae166ef95f9bff989d964061cbf81a0c",
      "signature": "811ec78f7fefcabbda4bfa93b3eb67d9ae166ef95f9bff989d964061cbf81a0c",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "811ec78f7fefcabbda4bfa93b3eb67d9ae166ef95f9bff989d964061cbf81a0c",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "bundled:///libs/lib.es2018.regexp.d.ts",
      "version": "717937616a17072082152a2ef351cb51f98802fb4b2fdabd32399843875974ca",
      "signature": "717937616a17072082152a2ef351cb51f98802fb4b2fdabd32399843875974ca",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "717937616a17072082152a2ef351cb51f98802fb4b2fdabd32399843875974ca",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "bundled:///libs/lib.es2019.array.d.ts",
      "version": "d7e7d9b7b50e5f22c915b525acc5a49a7a6584cf8f62d0569e557c5cfc4b2ac2",
      "signature": "d7e7d9b7b50e5f22c915b525acc5a49a7a6584cf8f62d0569e557c5cfc4b2ac2",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "d7e7d9b7b50e5f22c915b525acc5a49a7a6584cf8f62d0569e557c5cfc4b2ac2",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "bundled:///libs/lib.es2019.object.d.ts",
      "version": "71c37f4c9543f31dfced6c7840e068c5a5aacb7b89111a4364b1d5276b852557",
      "signature": "71c37f4c9543f31dfced6c7840e068c5a5aacb7b89111a4364b1d5276b852557",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "71c37f4c9543f31dfced6c7840e068c5a5aacb7b89111a4364b1d5276b852557",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "bundled:///libs/lib.es2019.string.d.ts",
      "version": "576711e016cf4f1804676043e6a0a5414252560eb57de9faceee34d79798c850",
      "signature": "576711e016cf4f1804676043e6a0a5414252560eb57de9faceee34d79798c850",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "576711e016cf4f1804676043e6a0a5414252560eb57de9faceee34d79798c850",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "bundled:///libs/lib.es2019.symbol.d.ts",
      "version": "89c1b1281ba7b8a96efc676b11b264de7a8374c5ea1e6617f11880a13fc56dc6",
      "signature": "89c1b1281ba7b8a96efc676b11b264de7a8374c5ea1e6617f11880a13fc56dc6",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "89c1b1281ba7b8a96efc676b11b264de7a8374c5ea1e6617f11880a13fc56dc6",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "bundled:///libs/lib.es2019.intl.d.ts",
      "version": "74f7fa2d027d5b33eb0471c8e82a6c87216223181ec31247c357a3e8e2fddc5b",
      "signature": "74f7fa2d027d5b33eb0471c8e82a6c87216223181ec31247c357a3e8e2fddc5b",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "74f7fa2d027d5b33eb0471c8e82a6c87216223181ec31247c357a3e8e2fddc5b",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "bundled:///libs/lib.es2020.bigint.d.ts",
      "version": "d6d7ae4d1f1f3772e2a3cde568ed08991a8ae34a080ff1151af28b7f798e22ca",
      "signature": "d6d7ae4d1f1f3772e2a3cde568ed08991a8ae34a080ff1151af28b7f798e22ca",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "d6d7ae4d1f1f3772e2a3cde568ed08991a8ae34a080ff1151af28b7f798e22ca",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "bundled:///libs/lib.es2020.date.d.ts",
      "version": "063600664504610fe3e99b717a1223f8b1900087fab0b4cad1496a114744f8df",
      "signature": "063600664504610fe3e99b717a1223f8b1900087fab0b4cad1496a114744f8df",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "063600664504610fe3e99b717a1223f8b1900087fab0b4cad1496a114744f8df",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "bundled:///libs/lib.es2020.promise.d.ts",
      "version": "934019d7e3c81950f9a8426d093458b65d5aff2c7c1511233c0fd5b941e608ab",
      "signature": "934019d7e3c81950f9a8426d093458b65d5aff2c7c1511233c0fd5b941e608ab",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "934019d7e3c81950f9a8426d093458b65d5aff2c7c1511233c0fd5b941e608ab",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "bundled:///libs/lib.es2020.sharedmemory.d.ts",
      "version": "52ada8e0b6e0482b728070b7639ee42e83a9b1c22d205992756fe020fd9f4a47",
      "signature": "52ada8e0b6e0482b728070b7639ee42e83a9b1c22d205992756fe020fd9f4a47",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "52ada8e0b6e0482b728070b7639ee42e83a9b1c22d205992756fe020fd9f4a47",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "bundled:///libs/lib.es2020.string.d.ts",
      "version": "3bdefe1bfd4d6dee0e26f928f93ccc128f1b64d5d501ff4a8cf3c6371200e5e6",
      "signature": "3bdefe1bfd4d6dee0e26f928f93ccc128f1b64d5d501ff4a8cf3c6371200e5e6",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "3bdefe1bfd4d6dee0e26f928f93ccc128f1b64d5d501ff4a8cf3c6371200e5e6",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "bundled:///libs/lib.es2020.symbol.wellknown.d.ts",
      "version": "59fb2c069260b4ba00b5643b907ef5d5341b167e7d1dbf58dfd895658bda2867",
      "signature": "59fb2c069260b4ba00b5643b907ef5d5341b167e7d1dbf58dfd895658bda2867",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "59fb2c069260b4ba00b5643b907ef5d5341b167e7d1dbf58dfd895658bda2867",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "bundled:///libs/lib.es2020.intl.d.ts",
      "version": "639e512c0dfc3fad96a84caad71b8834d66329a1f28dc95e3946c9b58176c73a",
      "signature": "639e512c0dfc3fad96a84caad71b8834d66329a1f28dc95e3946c9b58176c73a",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "639e512c0dfc3fad96a84caad71b8834d66329a1f28dc95e3946c9b58176c73a",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "bundled:///libs/lib.es2020.number.d.ts",
      "version": "368af93f74c9c932edd84c58883e736c9e3d53cec1fe24c0b0ff451f529ceab1",
      "signature": "368af93f74c9c932edd84c58883e736c9e3d53cec1fe24c0b0ff451f529ceab1",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "368af93f74c9c932edd84c58883e736c9e3d53cec1fe24c0b0ff451f529ceab1",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "bundled:///libs/lib.es2021.promise.d.ts",
      "version": "af3dd424cf267428f30ccfc376f47a2c0114546b55c44d8c0f1d57d841e28d74",
      "signature": "af3dd424cf267428f30ccfc376f47a2c0114546b55c44d8c0f1d57d841e28d74",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "af3dd424cf267428f30ccfc376f47a2c0114546b55c44d8c0f1d57d841e28d74",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "bundled:///libs/lib.es2021.string.d.ts",
      "version": "995c005ab91a498455ea8dfb63aa9f83fa2ea793c3d8aa344be4a1678d06d399",
      "signature": "995c005ab91a498455ea8dfb63aa9f83fa2ea793c3d8aa344be4a1678d06d399",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "995c005ab91a498455ea8dfb63aa9f83fa2ea793c3d8aa344be4a1678d06d399",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "bundled:///libs/lib.es2021.weakref.d.ts",
      "version": "959d36cddf5e7d572a65045b876f2956c973a586da58e5d26cde519184fd9b8a",
      "signature": "959d36cddf5e7d572a65045b876f2956c973a586da58e5d26cde519184fd9b8a",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "959d36cddf5e7d572a65045b876f2956c973a586da58e5d26cde519184fd9b8a",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "bundled:///libs/lib.es2021.intl.d.ts",
      "version": "965f36eae237dd74e6cca203a43e9ca801ce38824ead814728a2807b1910117d",
      "signature": "965f36eae237dd74e6cca203a43e9ca801ce38824ead814728a2807b1910117d",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "965f36eae237dd74e6cca203a43e9ca801ce38824ead814728a2807b1910117d",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "bundled:///libs/lib.es2022.array.d.ts",
      "version": "3925a6c820dcb1a06506c90b1577db1fdbf7705d65b62b99dce4be75c637e26b",
      "signature": "3925a6c820dcb1a06506c90b1577db1fdbf7705d65b62b99dce4be75c637e26b",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "3925a6c820dcb1a06506c90b1577db1fdbf7705d65b62b99dce4be75c637e26b",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "bundled:///libs/lib.es2022.error.d.ts",
      "version": "0a3d63ef2b853447ec4f749d3f368ce642264246e02911fcb1590d8c161b8005",
      "signature": "0a3d63ef2b853447ec4f749d3f368ce642264246e02911fcb1590d8c161b8005",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "0a3d63ef2b853447ec4f749d3f368ce642264246e02911fcb1590d8c161b8005",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "bundled:///libs/lib.es2022.intl.d.ts",
      "version": "b5ce7a470bc3628408429040c4e3a53a27755022a32fd05e2cb694e7015386c7",
      "signature": "b5ce7a470bc3628408429040c4e3a53a27755022a32fd05e2cb694e7015386c7",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "b5ce7a470bc3628408429040c4e3a53a27755022a32fd05e2cb694e7015386c7",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "bundled:///libs/lib.es2022.object.d.ts",
      "version": "8444af78980e3b20b49324f4a16ba35024fef3ee069a0eb67616ea6ca821c47a",
      "signature": "8444af78980e3b20b49324f4a16ba35024fef3ee069a0eb67616ea6ca821c47a",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "8444af78980e3b20b49324f4a16ba35024fef3ee069a0eb67616ea6ca821c47a",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "bundled:///libs/lib.es2022.string.d.ts",
      "version": "3287d9d085fbd618c3971944b65b4be57859f5415f495b33a6adc994edd2f004",
      "signature": "3287d9d085fbd618c3971944b65b4be57859f5415f495b33a6adc994edd2f004",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "3287d9d085fbd618c3971944b65b4be57859f5415f495b33a6adc994edd2f004",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "bundled:///libs/lib.es2022.regexp.d.ts",
      "version": "b4b67b1a91182421f5df999988c690f14d813b9850b40acd06ed44691f6727ad",
      "signature": "b4b67b1a91182421f5df999988c690f14d813b9850b40acd06ed44691f6727ad",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "b4b67b1a91182421f5df999988c690f14d813b9850b40acd06ed44691f6727ad",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "bundled:///libs/lib.es2023.array.d.ts",
      "version": "df83c2a6c73228b625b0beb6669c7ee2a09c914637e2d35170723ad49c0f5cd4",
      "signature": "df83c2a6c73228b625b0beb6669c7ee2a09c914637e2d35170723ad49c0f5cd4",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "df83c2a6c73228b625b0beb6669c7ee2a09c914637e2d35170723ad49c0f5cd4",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "bundled:///libs/lib.es2023.collection.d.ts",
      "version": "436aaf437562f276ec2ddbee2f2cdedac7664c1e4c1d2c36839ddd582eeb3d0a",
      "signature": "436aaf437562f276ec2ddbee2f2cdedac7664c1e4c1d2c36839ddd582eeb3d0a",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "436aaf437562f276ec2ddbee2f2cdedac7664c1e4c1d2c36839ddd582eeb3d0a",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "bundled:///libs/lib.es2023.intl.d.ts",
      "version": "8e3c06ea092138bf9fa5e874a1fdbc9d54805d074bee1de31b99a11e2fec239d",
      "signature": "8e3c06ea092138bf9fa5e874a1fdbc9d54805d074bee1de31b99a11e2fec239d",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "8e3c06ea092138bf9fa5e874a1fdbc9d54805d074bee1de31b99a11e2fec239d",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "bundled:///libs/lib.es2024.arraybuffer.d.ts",
      "version": "87dc0f382502f5bbce5129bdc0aea21e19a3abbc19259e0b43ae038a9fc4e326",
      "signature": "87dc0f382502f5bbce5129bdc0aea21e19a3abbc19259e0b43ae038a9fc4e326",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "87dc0f382502f5bbce5129bdc0aea21e19a3abbc19259e0b43ae038a9fc4e326",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "bundled:///libs/lib.es2024.collection.d.ts",
      "version": "b1cb28af0c891c8c96b2d6b7be76bd394fddcfdb4709a20ba05a7c1605eea0f9",
      "signature": "b1cb28af0c891c8c96b2d6b7be76bd394fddcfdb4709a20ba05a7c1605eea0f9",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "b1cb28af0c891c8c96b2d6b7be76bd394fddcfdb4709a20ba05a7c1605eea0f9",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "bundled:///libs/lib.es2024.object.d.ts",
      "version": "2fef54945a13095fdb9b84f705f2b5994597640c46afeb2ce78352fab4cb3279",
      "signature": "2fef54945a13095fdb9b84f705f2b5994597640c46afeb2ce78352fab4cb3279",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "2fef54945a13095fdb9b84f705f2b5994597640c46afeb2ce78352fab4cb3279",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "bundled:///libs/lib.es2024.promise.d.ts",
      "version": "ac77cb3e8c6d3565793eb90a8373ee8033146315a3dbead3bde8db5eaf5e5ec6",
      "signature": "ac77cb3e8c6d3565793eb90a8373ee8033146315a3dbead3bde8db5eaf5e5ec6",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "ac77cb3e8c6d3565793eb90a8373ee8033146315a3dbead3bde8db5eaf5e5ec6",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "bundled:///libs/lib.es2024.regexp.d.ts",
      "version": "56e4ed5aab5f5920980066a9409bfaf53e6d21d3f8d020c17e4de584d29600ad",
      "signature": "56e4ed5aab5f5920980066a9409bfaf53e6d21d3f8d020c17e4de584d29600ad",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "56e4ed5aab5f5920980066a9409bfaf53e6d21d3f8d020c17e4de584d29600ad",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "bundled:///libs/lib.es2024.sharedmemory.d.ts",
      "version": "4ece9f17b3866cc077099c73f4983bddbcb1dc7ddb943227f1ec070f529dedd1",
      "signature": "4ece9f17b3866cc077099c73f4983bddbcb1dc7ddb943227f1ec070f529dedd1",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "4ece9f17b3866cc077099c73f4983bddbcb1dc7ddb943227f1ec070f529dedd1",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "bundled:///libs/lib.es2024.string.d.ts",
      "version": "0a6282c8827e4b9a95f4bf4f5c205673ada31b982f50572d27103df8ceb8013c",
      "signature": "0a6282c8827e4b9a95f4bf4f5c205673ada31b982f50572d27103df8ceb8013c",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "0a6282c8827e4b9a95f4bf4f5c205673ada31b982f50572d27103df8ceb8013c",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "bundled:///libs/lib.esnext.array.d.ts",
      "version": "1c9319a09485199c1f7b0498f2988d6d2249793ef67edda49d1e584746be9032",
      "signature": "1c9319a09485199c1f7b0498f2988d6d2249793ef67edda49d1e584746be9032",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "1c9319a09485199c1f7b0498f2988d6d2249793ef67edda49d1e584746be9032",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "bundled:///libs/lib.esnext.collection.d.ts",
      "version": "e3a2a0cee0f03ffdde24d89660eba2685bfbdeae955a6c67e8c4c9fd28928eeb",
      "signature": "e3a2a0cee0f03ffdde24d89660eba2685bfbdeae955a6c67e8c4c9fd28928eeb",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "e3a2a0cee0f03ffdde24d89660eba2685bfbdeae955a6c67e8c4c9fd28928eeb",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "bundled:///libs/lib.esnext.intl.d.ts",
      "version": "811c71eee4aa0ac5f7adf713323a5c41b0cf6c4e17367a34fbce379e12bbf0a4",
      "signature": "811c71eee4aa0ac5f7adf713323a5c41b0cf6c4e17367a34fbce379e12bbf0a4",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "811c71eee4aa0ac5f7adf713323a5c41b0cf6c4e17367a34fbce379e12bbf0a4",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "bundled:///libs/lib.esnext.disposable.d.ts",
      "version": "51ad4c928303041605b4d7ae32e0c1ee387d43a24cd6f1ebf4a2699e1076d4fa",
      "signature": "51ad4c928303041605b4d7ae32e0c1ee387d43a24cd6f1ebf4a2699e1076d4fa",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "51ad4c928303041605b4d7ae32e0c1ee387d43a24cd6f1ebf4a2699e1076d4fa",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "bundled:///libs/lib.esnext.promise.d.ts",
      "version": "60037901da1a425516449b9a20073aa03386cce92f7a1fd902d7602be3a7c2e9",
      "signature": "60037901da1a425516449b9a20073aa03386cce92f7a1fd902d7602be3a7c2e9",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "60037901da1a425516449b9a20073aa03386cce92f7a1fd902d7602be3a7c2e9",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "bundled:///libs/lib.esnext.decorators.d.ts",
      "version": "d4b1d2c51d058fc21ec2629fff7a76249dec2e36e12960ea056e3ef89174080f",
      "signature": "d4b1d2c51d058fc21ec2629fff7a76249dec2e36e12960ea056e3ef89174080f",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "d4b1d2c51d058fc21ec2629fff7a76249dec2e36e12960ea056e3ef89174080f",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "bundled:///libs/lib.esnext.iterator.d.ts",
      "version": "22adec94ef7047a6c9d1af3cb96be87a335908bf9ef386ae9fd50eeb37f44c47",
      "signature": "22adec94ef7047a6c9d1af3cb96be87a335908bf9ef386ae9fd50eeb37f44c47",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "22adec94ef7047a6c9d1af3cb96be87a335908bf9ef386ae9fd50eeb37f44c47",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "bundled:///libs/lib.esnext.float16.d.ts",
      "version": "4245fee526a7d1754529d19227ecbf3be066ff79ebb6a380d78e41648f2f224d",
      "signature": "4245fee526a7d1754529d19227ecbf3be066ff79ebb6a380d78e41648f2f224d",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "4245fee526a7d1754529d19227ecbf3be066ff79ebb6a380d78e41648f2f224d",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "bundled:///libs/lib.decorators.d.ts",
      "version": "8e7f8264d0fb4c5339605a15daadb037bf238c10b654bb3eee14208f860a32ea",
      "signature": "8e7f8264d0fb4c5339605a15daadb037bf238c10b654bb3eee14208f860a32ea",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "8e7f8264d0fb4c5339605a15daadb037bf238c10b654bb3eee14208f860a32ea",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "bundled:///libs/lib.decorators.legacy.d.ts",
      "version": "782dec38049b92d4e85c1585fbea5474a219c6984a35b004963b00beb1aab538",
      "signature": "782dec38049b92d4e85c1585fbea5474a219c6984a35b004963b00beb1aab538",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "782dec38049b92d4e85c1585fbea5474a219c6984a35b004963b00beb1aab538",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "bundled:///libs/lib.esnext.full.d.ts",
      "version": "bde31fd423cd93b0eff97197a3f66df7c93e8c0c335cbeb113b7ff1ac35c23f4",
      "signature": "bde31fd423cd93b0eff97197a3f66df7c93e8c0c335cbeb113b7ff1ac35c23f4",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "../compiler/parser.d.ts",
      "version": "2e29cd9a98755c46896f7a2d56524db2d6d96b248e36db46de14c30bf47c8d05",
      "signature": "2e29cd9a98755c46896f7a2d56524db2d6d96b248e36db46de14c30bf47c8d05",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "../../src/services/services.ts",
      "version": "407537635fda1a543a422ecdd456c1402aaa2083cde5acfb4eb424ab02fc0612",
      "signature": "8e609bb71c20b858c77f0e9f90bb1319db8477b13f9f965f1a1e18524bf50881",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "407537635fda1a543a422ecdd456c1402aaa2083cde5acfb4eb424ab02fc0612",
        "signature": "8e609bb71c20b858c77f0e9f90bb1319db8477b13f9f965f1a1e18524bf50881",
        "impliedNodeFormat": 1
      }
    }
  ],
  "fileIdsList": [
    [
      "../compiler/parser.d.ts"
    ]
  ],
  "options": {
    "composite": true,
    "module": 199,
    "outDir": "./",
    "rewriteRelativeImportExtensions": true,
    "rootDir": "../../src/services"
  },
  "referencedMap": {
    "../../src/services/services.ts": [
      "../compiler/parser.d.ts"
    ]
  },
  "latestChangedDtsFile": "./services.d.ts",
  "size": 13963
}
//// [/home/src/workspaces/solution/src/compiler/parser.ts] no change
//// [/home/src/workspaces/solution/src/compiler/tsconfig.json] no change
//// [/home/src/workspaces/solution/src/services/services.ts] no change
//// [/home/src/workspaces/solution/src/services/tsconfig.json] no change
//// [/home/src/workspaces/solution/src/tsconfig-base.json] no change

