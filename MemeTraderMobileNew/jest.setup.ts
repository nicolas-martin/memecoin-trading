import '@testing-library/jest-native/extend-expect';

declare global {
  namespace jest {
    interface Matchers<R> {
      toHaveTextContent: (text: string) => R;
      toBeVisible: () => R;
      toBeEnabled: () => R;
      toBeDisabled: () => R;
    }
  }
} 