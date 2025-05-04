import { hydrateRoot } from 'react-dom/client';
import createApp from './createApp';

/*
  Index.html has the `/mount.js` in it's script tag.
  Each time an app is updated, index is refreshed from the server and recalls
  this file.
 */

// eslint-disable-next-line no-undef
hydrateRoot(document.querySelector('main'), createApp());
