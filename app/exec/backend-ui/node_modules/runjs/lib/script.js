'use strict';

Object.defineProperty(exports, "__esModule", {
  value: true
});
exports.requirer = requirer;
exports.hasAccess = hasAccess;
exports.getConfig = getConfig;
exports.load = load;
exports.describe = describe;
exports.call = call;
exports.main = main;

var _path = require('path');

var _path2 = _interopRequireDefault(_path);

var _fs = require('fs');

var _fs2 = _interopRequireDefault(_fs);

var _chalk = require('chalk');

var _chalk2 = _interopRequireDefault(_chalk);

var _lodash = require('lodash.padend');

var _lodash2 = _interopRequireDefault(_lodash);

var _microcli = require('microcli');

var _microcli2 = _interopRequireDefault(_microcli);

var _omelette = require('omelette');

var _omelette2 = _interopRequireDefault(_omelette);

var _common = require('./common');

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

const DEFAULT_RUNFILE_PATH = './runfile.js';
function requirer(filePath) {
  return require(_path2.default.resolve(filePath));
}

function hasAccess(filePath) {
  return _fs2.default.accessSync(_path2.default.resolve(filePath));
}

function getConfig(filePath) {
  let config;
  try {
    config = requirer(filePath).runjs || {};
  } catch (error) {
    config = {};
  }
  return config;
}

function load(config, logger, requirer, access) {
  const runfilePath = config['runfile'] || DEFAULT_RUNFILE_PATH;
  // Load requires if given in config
  if (Array.isArray(config['requires'])) {
    config['requires'].forEach(modulePath => {
      logger.log(_chalk2.default.gray(`Requiring ${modulePath}...`));
      requirer(modulePath);
    });
  }

  // Process runfile
  logger.log(_chalk2.default.gray(`Processing ${runfilePath}...`));

  try {
    access(runfilePath);
  } catch (error) {
    throw new _common.RunJSError(`No ${runfilePath} defined in ${process.cwd()}`);
  }

  const runfile = requirer(runfilePath);
  if (runfile.default) {
    return runfile.default;
  }
  return runfile;
}

function describe(obj, logger, namespace) {
  if (!namespace) {
    logger.log(_chalk2.default.yellow('Available tasks:'));
  }

  Object.keys(obj).forEach(key => {
    const value = obj[key];
    const nextNamespace = namespace ? `${namespace}:${key}` : key;
    const help = value.help;

    if (typeof value === 'function') {
      // Add task name
      const funcParams = help && help.params;
      let logArgs = [_chalk2.default.bold(nextNamespace)];

      // Add task params
      if (Array.isArray(funcParams) && funcParams.length) {
        logArgs[0] += ` [${funcParams.join(' ')}]`;
      }

      // Add description
      if (help && (help.description || typeof help === 'string')) {
        const description = help.description || help;
        logArgs[0] = (0, _lodash2.default)(logArgs[0], 40); // format
        logArgs.push('-', description.split('\n')[0]);
      }

      // Log
      logger.log(...logArgs);
    } else if (typeof value === 'object') {
      describe(value, logger, nextNamespace);
    }
  });

  if (!namespace) {
    logger.log('\n' + _chalk2.default.blue('Type "run [taskname] --help" to get more info if available.'));
  }
}

function tasks(obj, namespace) {
  let list = [];
  Object.keys(obj).forEach(key => {
    const value = obj[key];
    const nextNamespace = namespace ? `${namespace}:${key}` : key;

    if (typeof value === 'function') {
      list.push(nextNamespace);
    } else if (typeof value === 'object') {
      list = list.concat(tasks(value, nextNamespace));
    }
  });
  return list;
}

function call(obj, args, logger, subtaskName) {
  const taskName = subtaskName || args[2];

  if (typeof obj[taskName] === 'function') {
    const cli = (0, _microcli2.default)(args.slice(1), obj[taskName].help, null, logger);

    cli((options, ...params) => {
      obj[taskName].apply({ options }, params);
    });
    return obj[taskName];
  }

  let namespaces = taskName.split(':');
  const rootNamespace = namespaces.shift();
  const nextSubtaskName = namespaces.join(':');

  if (obj[rootNamespace]) {
    const calledTask = call(obj[rootNamespace], args, logger, nextSubtaskName);
    if (calledTask) {
      return calledTask;
    }
  }

  if (!subtaskName) {
    throw new _common.RunJSError(`Task ${taskName} not found`);
  }
}

function autocomplete(config) {
  const logger = new _common.SilentLogger();
  const completion = (0, _omelette2.default)('run <task>');
  completion.on('task', ({ reply }) => {
    const runfile = load(config, logger, requirer, hasAccess);
    reply(tasks(runfile));
  });
  completion.init();
}

function main() {
  try {
    const config = getConfig('./package.json');
    autocomplete(config);
    const runfile = load(config, _common.logger, requirer, hasAccess);
    const ARGV = process.argv.slice();

    if (ARGV.length > 2) {
      call(runfile, ARGV, _common.logger);
    } else {
      describe(runfile, _common.logger);
    }
  } catch (error) {
    if (error instanceof _common.RunJSError || error instanceof _microcli.CLIError) {
      _common.logger.error(error.message);
      process.exit(1);
    } else {
      _common.logger.log(error);
      process.exit(1);
    }
  }
}