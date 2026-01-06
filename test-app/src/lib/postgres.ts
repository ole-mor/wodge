import { apiPost } from './wodge';

export interface QueryResult<T = any> {
  [key: string]: any;
}

export const postgres = {
  /**
   * Execute a SELECT query
   */
  async query<T = any>(query: string, args: any[] = []): Promise<T[]> {
    return apiPost('/postgres/query', { query, args });
  },

  /**
   * Execute an INSERT/UPDATE/DELETE query
   */
  async execute(query: string, args: any[] = []): Promise<{ rows_affected: number }> {
    return apiPost('/postgres/execute', { query, args });
  }
};
