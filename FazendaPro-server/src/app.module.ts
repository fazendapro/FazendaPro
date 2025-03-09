import { Module } from '@nestjs/common';
import { UsersModule } from './controller/users/users.module';
import { AppService } from './app.service';
import { AppController } from './app.controller';
import { DatabaseModule } from './database/database.module';
@Module({
  imports: [DatabaseModule, UsersModule],
  controllers: [AppController],
  providers: [AppService],
})
export class AppModule {}
