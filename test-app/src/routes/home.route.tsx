import React, { useEffect, useState } from 'react';
import { HealthService } from '../api/health';

export default function Home() {
  const [health, setHealth] = useState<{ status: string } | null>(null);
  useEffect(() => {
    HealthService.check().then(setHealth);
  }, []);
  if (!health) return <div>Loading...</div>;
  return (
    <div>
      <h1>Hello from Wodge!</h1>
      <p>This is a home route.</p>
      <p>Health: {health.status}</p>
    </div>
  );
}