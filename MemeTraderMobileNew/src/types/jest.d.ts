/// <reference types="jest" />

declare namespace jest {
  interface Expect {
    <T = any>(actual: T): jest.Matchers<T>;
  }

  interface InverseAsymmetricMatchers {
    toBeTruthy(): void;
    toBeFalsy(): void;
    toBeNull(): void;
  }

  interface Matchers<R, T = {}> {
    toBeTruthy(): R;
    toBeFalsy(): R;
    toBeNull(): R;
    toBeUndefined(): R;
    toBeDefined(): R;
    toBeNaN(): R;
    toEqual(expected: T): R;
    toHaveLength(length: number): R;
    toContain(item: any): R;
    toBeCalledWith(...args: any[]): R;
    toHaveBeenCalledWith(...args: any[]): R;
  }
}

declare module '@testing-library/react-native' {
  export interface RenderResult {
    getByText(text: string | RegExp): any;
    queryByText(text: string | RegExp): any | null;
    getAllByText(text: string | RegExp): any[];
    getByTestId(testId: string): any;
    queryByTestId(testId: string): any | null;
    getAllByTestId(testId: string): any[];
  }
}

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