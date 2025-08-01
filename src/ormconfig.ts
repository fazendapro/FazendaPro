import { DataSource } from 'typeorm';
import { User } from './features/user-management/users/users.entity';

export const AppDataSource = new DataSource({
  type: 'postgres',
  url: process.env.DATABASE_URL || 'postgresql://user:123456@postgres:5432/fazendapro_db',
  entities: [User],
  synchronize: process.env.NODE_ENV !== 'production',
  migrations: ['src/migrations/*.ts'],
  migrationsTableName: 'migrations',
});
