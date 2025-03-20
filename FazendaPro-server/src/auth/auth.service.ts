import { Injectable, UnauthorizedException } from '@nestjs/common';
import { UsersService } from '../features/user-management/users/users.service';
import { JwtService } from '@nestjs/jwt';

@Injectable()
export class AuthService {
  constructor(
    private usersService: UsersService,
    private jwtService: JwtService,
  ) {}

  async signIn(
    username: string,
    pass: string,
  ): Promise<{ access_token: string }> {
    const user = await this.usersService.findOneByUsername(username);

    if (!user) {
      console.error('Usuário não encontrado');
      throw new UnauthorizedException();
    }

    if (user.password !== pass) {
      console.error('Senha incorreta');
      throw new UnauthorizedException();
    }

    const payload = { sub: user.id, username: user.firstName };
    return {
      access_token: await this.jwtService.signAsync(payload),
    };
  }
}
