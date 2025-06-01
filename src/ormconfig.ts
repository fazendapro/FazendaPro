import { DataSource } from 'typeorm';
import { User } from './features/user-management/users/users.entity';

export const AppDataSource = new DataSource({
  type: 'mysql',
  url: process.env.JAWSDB_URL || 'mysql://user:123456@localhost:3306/fazendapro_db',
  entities: [User],
  synchronize: process.env.NODE_ENV !== 'production',
  migrations: ['src/migrations/*.ts'],
  migrationsTableName: 'migrations',
});
