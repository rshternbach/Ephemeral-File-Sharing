import log from 'loglevel';
import { config } from '../config';

const logLevel = config.MODE === 'production' ? 'warn' : 'debug';
log.setLevel(logLevel);

const originalFactory = log.methodFactory;
log.methodFactory = function (methodName, logLevel, loggerName) {
  const rawMethod = originalFactory(methodName, logLevel, loggerName);
  return function (message, ...args) {
    rawMethod(`[${new Date().toISOString()}] [${methodName.toUpperCase()}] ${message}`, ...args);
  };
};

log.setLevel(log.getLevel());

export default log;