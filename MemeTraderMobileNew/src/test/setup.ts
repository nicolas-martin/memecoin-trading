import '@testing-library/jest-native/extend-expect';
import { configure } from '@testing-library/react-native';

// Configure testing library
configure({
  asyncUtilTimeout: 5000,
  defaultHidden: true,
  testIdAttribute: 'testID',
});

// Mock the Image component
jest.mock('react-native/Libraries/Image/Image', () => ({
  __esModule: true,
  default: 'Image',
}));

// Mock require for placeholder images
jest.mock('../assets/placeholder.png', () => 'placeholder-image'); 