// User Entity - Corresponds to users table (RF01)
// Example TypeORM entity for NestJS implementation

import {
  Entity,
  PrimaryGeneratedColumn,
  Column,
  CreateDateColumn,
  UpdateDateColumn,
  DeleteDateColumn,
  OneToMany,
  Index,
} from 'typeorm';
import { Animal } from './animal.entity';
import { Vaccine } from './vaccine.entity';
import { Supplier } from './supplier.entity';

@Entity('users')
@Index(['email'])
export class User {
  @PrimaryGeneratedColumn()
  id: number;

  @Column({ type: 'varchar', length: 255 })
  name: string;

  @Column({ type: 'varchar', length: 255, unique: true })
  @Index()
  email: string;

  @Column({ type: 'varchar', length: 255, name: 'password_hash' })
  passwordHash: string;

  @CreateDateColumn({ name: 'created_at' })
  createdAt: Date;

  @UpdateDateColumn({ name: 'updated_at' })
  updatedAt: Date;

  @DeleteDateColumn({ name: 'deleted_at' })
  deletedAt?: Date;

  // Relationships
  @OneToMany(() => Animal, animal => animal.user)
  animals: Animal[];

  @OneToMany(() => Vaccine, vaccine => vaccine.user)
  vaccines: Vaccine[];

  @OneToMany(() => Supplier, supplier => supplier.user)
  suppliers: Supplier[];
}