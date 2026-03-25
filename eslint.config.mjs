import nextCoreWebVitals from 'eslint-config-next/core-web-vitals'
import tseslint from 'typescript-eslint'

export default [
    ...nextCoreWebVitals,
    ...tseslint.configs.strict,
    {
        rules: {
            '@typescript-eslint/no-unused-vars': ['error', { argsIgnorePattern: '^_', varsIgnorePattern: '^_' }],
            '@typescript-eslint/no-explicit-any': 'error',
            '@typescript-eslint/consistent-type-imports': 'error',
            'no-console': 'warn',
            'prefer-const': 'error',
            'no-var': 'error',
            eqeqeq: ['error', 'always'],
            curly: ['error', 'all'],
            'no-throw-literal': 'error',
            'no-eval': 'error',
            'no-implied-eval': 'error',
            'no-new-wrappers': 'error',
            'no-restricted-imports': 'off',
        },
    },
    {
        ignores: ['.next/', 'node_modules/', 'dist/', 'build/'],
    },
]
