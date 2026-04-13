module.exports = {
  root: true,
  env: {
    browser: true,
    es2022: true,
    node: true
  },
  extends: [
    'eslint:recommended',
    'plugin:vue/vue3-recommended'
  ],
  parserOptions: {
    ecmaVersion: 'latest',
    sourceType: 'module'
  },
  rules: {
    // Template errors (catches invalid tags, unclosed elements)
    'vue/no-parsing-error': 'error',
    'vue/valid-template-root': 'error',
    'vue/no-duplicate-attributes': 'error',
    'vue/no-unused-vars': 'warn',
    'vue/require-v-for-key': 'error',

    // Relax some rules for rapid development
    'vue/multi-word-component-names': 'off',
    'vue/max-attributes-per-line': 'off',
    'vue/singleline-html-element-content-newline': 'off',
    'vue/html-self-closing': 'off',
    'vue/html-closing-bracket-newline': 'off',
    'vue/first-attribute-linebreak': 'off',
    'vue/html-indent': 'off',
    'vue/attribute-hyphenation': 'off',
    'vue/attributes-order': 'off',

    // JS
    'no-unused-vars': ['warn', { argsIgnorePattern: '^_' }],
    'no-console': 'warn'
  }
}
