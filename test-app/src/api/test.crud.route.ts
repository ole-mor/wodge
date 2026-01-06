import { postgres } from '@/lib/postgres';

export const TestService = {
  async list() {
    return postgres.query('SELECT * FROM test');
  },

  async get(id: string) {
    const rows = await postgres.query('SELECT * FROM test WHERE id = $1', [id]);
    return rows[0] || null;
  },

  async create(data: Record<string, any>) {
    const keys = Object.keys(data);
    const values = Object.values(data);
    const placeholders = keys.map((_, i) => '$' + (i + 1)).join(', ');
    const columns = keys.join(', ');
    
    const query = 'INSERT INTO test (' + columns + ') VALUES (' + placeholders + ')';
    return postgres.execute(query, values);
  },

  async delete(id: string) {
    return postgres.execute('DELETE FROM test WHERE id = $1', [id]);
  }
};
