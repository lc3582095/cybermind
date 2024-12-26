import React from 'react';
import { createBrowserRouter } from 'react-router-dom';
import Layout from '../layouts/MainLayout';
import Chat from '../pages/Chat';
import Image from '../pages/Image';
import Video from '../pages/Video';
import Music from '../pages/Music';
import PPT from '../pages/PPT';
import Apps from '../pages/Apps';

const router = createBrowserRouter([
  {
    path: '/',
    element: <Layout />,
    children: [
      {
        path: '/',
        element: <Chat />,
      },
      {
        path: '/chat',
        element: <Chat />,
      },
      {
        path: '/image',
        element: <Image />,
      },
      {
        path: '/video',
        element: <Video />,
      },
      {
        path: '/music',
        element: <Music />,
      },
      {
        path: '/ppt',
        element: <PPT />,
      },
      {
        path: '/apps',
        element: <Apps />,
      },
    ],
  },
]);

export default router; 