import React from 'react';
import { useRoutes } from 'react-router-dom';
import { routes } from './routes.generated';

export default function App() {
  return useRoutes(routes);
}