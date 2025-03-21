import { Module } from '@nestjs/common';
import { TypeOrmModule } from '@nestjs/typeorm';
import { User } from '../features/user-management/users/users.entity';

@Module({
  imports: [
    TypeOrmModule.forRoot({
      type: 'mysql',
      host: 'mysql',
      port: 3306,
      username: 'user',
      password: '123456',
      database: 'fazendapro_db',
      entities: [User],
      synchronize: true,
      autoLoadEntities: true,
    }),
  ],
})

// @Module({
//   imports: [
//     TypeOrmModule.forRoot({
//       type: 'mysql',
//       url: process.env.JAWSDB_URL || 'mysql://user:123456@mysql:3306/fazendapro_db',
//       entities: [User],
//       synchronize: process.env.NODE_ENV !== 'production',
//       autoLoadEntities: true,
//     }),
//   ],
// })

export class DatabaseModule {}