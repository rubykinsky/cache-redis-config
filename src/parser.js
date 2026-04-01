// parser.js
const fs = require('fs');
const path = require('path');

/**
 * Parses a configuration file (JSON or JavaScript) and returns the configuration object.
 *
 * @param {string} filePath The path to the configuration file.
 * @returns {object} The configuration object.
 * @throws {Error} If the file does not exist or if there is an error parsing the file.
 */
function parseConfigFile(filePath) {
  if (!fs.existsSync(filePath)) {
    throw new Error(`Configuration file not found: ${filePath}`);
  }

  const fileExtension = path.extname(filePath).toLowerCase();

  try {
    switch (fileExtension) {
      case '.json':
        const rawData = fs.readFileSync(filePath, 'utf8');
        return JSON.parse(rawData);
      case '.js':
        // eslint-disable-next-line import/no-dynamic-require, global-require
        return require(filePath);
      default:
        throw new Error(`Unsupported file type: ${fileExtension}. Supported types are .json and .js.`);
    }
  } catch (error) {
    throw new Error(`Error parsing configuration file: ${filePath}. ${error.message}`);
  }
}

module.exports = {
  parseConfigFile,
};