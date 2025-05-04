import App from './App.jsx';
/*
  `createApp` is for anything that does not need the window object or document, specifically
  as it's not loaded from the server yet.
*/
export default function createApp() {
  return <App />;
}
