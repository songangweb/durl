'use strict';

Object.defineProperty(exports, "__esModule", {
  value: true
});
exports.logger = exports.SilentLogger = exports.Logger = exports.RunJSError = undefined;

var _chalk = require('chalk');

var _chalk2 = _interopRequireDefault(_chalk);

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

// Needed to use ES5 inheritance, because of issues with Error subclassing for Babel
class RunJSError extends Error {
  constructor(message) {
    message = message && message.split('\n')[0]; // assign only first line
    super(message);
  }
}

exports.RunJSError = RunJSError;
class Logger {
  title(...args) {
    console.log(_chalk2.default.bold(...args));
  }
  log(...args) {
    console.log(...args);
  }
  warning(...args) {
    console.warn(_chalk2.default.yellow(...args));
  }
  error(...args) {
    console.error(_chalk2.default.red(...args));
  }
}

exports.Logger = Logger;
class SilentLogger {
  title() {}
  log() {}
  warning() {}
  error() {}
}

exports.SilentLogger = SilentLogger;
const logger = exports.logger = new Logger();