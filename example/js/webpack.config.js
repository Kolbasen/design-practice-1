const path = require('path');
const WebpackObfuscator = require('webpack-obfuscator');

module.exports = (env = {}) => {

    console.log(env)

    const { ENTRY, SHOULD_OBFUSCATE, FILENAME = 'index.js' } = env;

    const obfuscatePlugin = Boolean(SHOULD_OBFUSCATE) ?
        [
            new WebpackObfuscator ({
                rotateStringArray: true
            }, [])
        ]
        : []

    console.log(obfuscatePlugin);

    return {
        mode: 'production',
        entry: ENTRY.split(','),
        output: {
            path: path.resolve(__dirname),
            filename: `${FILENAME}.js`,
          },
        plugins: obfuscatePlugin,
    }
}