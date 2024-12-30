# Meme Trader Mobile App

A React Native mobile application for tracking and analyzing meme coins. Built with Expo and TypeScript.

## Features

- Real-time meme coin price tracking
- Market cap and volume analytics
- 24-hour price change indicators
- Clean, modern UI with responsive design
- Integration with Go backend services

## Tech Stack

- React Native with Expo
- TypeScript
- React Navigation for routing
- Axios for API calls
- React Native Chart Kit for price charts
- Custom components for coin listings

## Getting Started

### Prerequisites

- Node.js (v16 or newer)
- npm or yarn
- Expo CLI
- iOS Simulator (for iOS development) or Android Studio (for Android development)
- Expo Go app (for physical device testing)

### Installation

1. Install dependencies:
```bash
make install-mobile-deps
```

2. Start the development server:
```bash
make start-mobile
```

### Running the App

You have three options to run the app:

1. Physical Device (Recommended):
   - Install Expo Go on your mobile device
   - Scan the QR code with:
     - iOS: Use the Camera app
     - Android: Use the Expo Go app
   - Make sure your device is on the same WiFi network as your computer

2. iOS Simulator:
   - Install Xcode
   - Run `make start-mobile-ios`

3. Android Emulator:
   - Install Android Studio
   - Set up an Android Virtual Device
   - Run `make start-mobile-android`

## Project Structure

```
src/
  ├── components/         # Reusable UI components
  │   └── CoinList/      # Coin listing components
  ├── screens/           # Screen components
  ├── services/          # API and other services
  ├── navigation/        # Navigation configuration
  └── assets/           # Images and other static assets
```

## Testing

Run the test suite:
```bash
npm test
```

## Development

- The app uses TypeScript for type safety
- Components are organized by feature
- API calls are centralized in the services directory
- Follows React Native best practices and conventions

## Backend Integration

The app integrates with a Go backend service that provides:
- Meme coin price data
- Market analytics
- Real-time updates

Make sure the backend service is running (`make start-backend`) before starting the app.
