const fs = require('fs');
const path = require('path');

class Parser {
  constructor(configFile) {
    this.configFile = configFile;
  }

  parse() {
    try {
      const rawData = fs.readFileSync(this.configFile, 'utf8');
      const config = JSON.parse(rawData);
      return config;
    } catch (error) {
      if (error.code === 'ENOENT') {
        throw new Error(`Config file not found: ${this.configFile}`);
      } else if (error instanceof SyntaxError) {
        throw new Error(`Invalid JSON in config file: ${this.configFile}`);
      } else {
        throw error;
      }
    }
  }

  validateConfig(config) {
    if (!config || typeof config !== 'object') {
      throw new Error('Invalid config object');
    }

    const requiredProperties = ['host', 'port', 'password'];
    requiredProperties.forEach(property => {
      if (!Object.prototype.hasOwnProperty.call(config, property)) {
        throw new Error(`Missing required property: ${property}`);
      }
    });

    return config;
  }

  parseConfig() {
    const config = this.parse();
    return this.validateConfig(config);
  }
}

module.exports = Parser;