import { Injectable, UnauthorizedException } from '@nestjs/common';
import { UsersService } from '../features/user-management/users/users.service';
import { JwtService } from '@nestjs/jwt';
import * as bcrypt from 'bcryptjs';

@Injectable()
export class AuthService {
  constructor(
    private usersService: UsersService,
    private jwtService: JwtService,
  ) {}

  async signIn(
    email: string,
    password: string,
  ): Promise<{ access_token: string }> {
    const user = await this.usersService.findOneByEmail(email);

    if (!user) {
      console.error('Usuário não encontrado');
      throw new UnauthorizedException();
    }

    if (await bcrypt.compare(password, user.password)) {
      const payload = { sub: user.id, email: user.email };
      return {
        access_token: await this.jwtService.signAsync(payload),
      };
    } else {
      console.error('Senha incorreta');
      throw new UnauthorizedException();
    }
  }
}
