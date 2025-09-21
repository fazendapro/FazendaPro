// TypeORM Database Configuration Example
// Configuration file for connecting to MySQL database

import { TypeOrmModuleOptions } from '@nestjs/typeorm';
import { ConfigService } from '@nestjs/config';

export const getDatabaseConfig = (configService: ConfigService): TypeOrmModuleOptions => ({
  type: 'mysql',
  host: configService.get<string>('DB_HOST', 'localhost'),
  port: configService.get<number>('DB_PORT', 3306),
  username: configService.get<string>('DB_USERNAME', 'root'),
  password: configService.get<string>('DB_PASSWORD', ''),
  database: configService.get<string>('DB_DATABASE', 'fazendapro'),
  charset: 'utf8mb4',
  timezone: 'Z',
  
  // Entity configuration
  entities: [
    __dirname + '/../**/*.entity{.ts,.js}',
  ],
  
  // Migration configuration
  migrations: [
    __dirname + '/../migrations/*{.ts,.js}',
  ],
  migrationsTableName: 'migrations',
  migrationsRun: false, // Set to true in production
  
  // Development settings
  synchronize: false, // Never use true in production
  dropSchema: false,
  logging: configService.get<string>('NODE_ENV') === 'development' ? ['query', 'error'] : ['error'],
  
  // Connection pool settings
  extra: {
    connectionLimit: 20,
    acquireTimeout: 60000,
    timeout: 60000,
  },
  
  // SSL configuration for production (like JawsDB on Heroku)
  ssl: configService.get<string>('NODE_ENV') === 'production' ? {
    rejectUnauthorized: false,
  } : false,
});

// Example environment variables (.env file):
/*
DB_HOST=localhost
DB_PORT=3306
DB_USERNAME=fazendapro_user
DB_PASSWORD=secure_password
DB_DATABASE=fazendapro

# For JawsDB on Heroku (example):
# DB_HOST=x1x1x1x1x1x1x1.us-east-1.rds.amazonaws.com
# DB_PORT=3306
# DB_USERNAME=username
# DB_PASSWORD=password
# DB_DATABASE=database_name
*/