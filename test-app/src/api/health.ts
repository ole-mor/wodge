import { apiGet } from '@/lib/wodge';

export const HealthService = {
  async check(): Promise<{ status: string }> {
    return apiGet('/health');
  }
};
