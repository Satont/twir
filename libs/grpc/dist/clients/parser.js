"use strict";
var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.createParser = void 0;
const nice_grpc_1 = require("nice-grpc");
const parser_js_1 = require("../generated/parser/parser.js");
const constants_js_1 = require("../servers/constants.js");
const helper_js_1 = require("./helper.js");
const createParser = (env) => __awaiter(void 0, void 0, void 0, function* () {
    const channel = (0, nice_grpc_1.createChannel)((0, helper_js_1.createClientAddr)(env, 'parser', constants_js_1.PORTS.PARSER_SERVER_PORT), nice_grpc_1.ChannelCredentials.createInsecure(), helper_js_1.CLIENT_OPTIONS);
    yield (0, helper_js_1.waitReady)(channel);
    const client = (0, nice_grpc_1.createClient)(parser_js_1.ParserDefinition, channel);
    return client;
});
exports.createParser = createParser;
//# sourceMappingURL=parser.js.map