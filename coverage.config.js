export default {
  provider: 'v8',
  reporter: ['text', 'json', 'html'],
  reportsDirectory: './coverage',
  exclude: [
    'node_modules/**',
    'dist/**',
    'coverage/**',
    '**/*.d.ts',
    '**/*.config.js',
    '**/*.config.ts',
    'src/test/**',
    'src/**/__tests__/**',
    'src/**/*.test.*',
    'src/**/*.spec.*',
  ],
  thresholds: {
    global: {
      branches: 70,
      functions: 70,
      lines: 70,
      statements: 70,
    },
  },
}
