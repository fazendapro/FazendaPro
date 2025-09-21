// Animal Entity - Corresponds to animals table (RF02, RF03)
// Core entity of the FazendaPro system

import {
  Entity,
  PrimaryGeneratedColumn,
  Column,
  CreateDateColumn,
  UpdateDateColumn,
  DeleteDateColumn,
  ManyToOne,
  OneToMany,
  JoinColumn,
  Index,
} from 'typeorm';
import { User } from './user.entity';
import { Lot } from './lot.entity';
import { AnimalVaccine } from './animal-vaccine.entity';
import { WeightRecord } from './weight-record.entity';
import { PregnancyRecord } from './pregnancy-record.entity';
import { Sale } from './sale.entity';

export enum AnimalGender {
  MALE = 'male',
  FEMALE = 'female',
}

export enum AnimalStatus {
  ACTIVE = 'active',
  SOLD = 'sold',
  DECEASED = 'deceased',
  TRANSFERRED = 'transferred',
}

@Entity('animals')
@Index(['ear_tag'], { unique: true })
@Index(['user_id'])
@Index(['status'])
@Index(['gender'])
export class Animal {
  @PrimaryGeneratedColumn()
  id: number;

  @Column({ name: 'user_id' })
  userId: number;

  @Column({ name: 'lot_id', nullable: true })
  lotId?: number;

  @Column({ type: 'varchar', length: 255 })
  name: string;

  @Column({ name: 'ear_tag', type: 'varchar', length: 100, unique: true })
  @Index()
  earTag: string;

  @Column({ name: 'birth_date', type: 'date' })
  birthDate: Date;

  @Column({ type: 'varchar', length: 255, nullable: true })
  race?: string;

  @Column({ 
    type: 'enum', 
    enum: AnimalGender,
    default: AnimalGender.FEMALE 
  })
  gender: AnimalGender;

  @Column({ name: 'mother_id', nullable: true })
  motherId?: number;

  @Column({ name: 'father_id', nullable: true })
  fatherId?: number;

  @Column({ 
    type: 'enum', 
    enum: AnimalStatus, 
    default: AnimalStatus.ACTIVE 
  })
  status: AnimalStatus;

  @Column({ 
    name: 'current_weight', 
    type: 'decimal', 
    precision: 8, 
    scale: 2, 
    nullable: true 
  })
  currentWeight?: number;

  @Column({ 
    name: 'daily_milk_production', 
    type: 'decimal', 
    precision: 8, 
    scale: 2, 
    default: 0 
  })
  dailyMilkProduction: number;

  @Column({ type: 'text', nullable: true })
  notes?: string;

  @CreateDateColumn({ name: 'created_at' })
  createdAt: Date;

  @UpdateDateColumn({ name: 'updated_at' })
  updatedAt: Date;

  @DeleteDateColumn({ name: 'deleted_at' })
  deletedAt?: Date;

  // Relationships
  @ManyToOne(() => User, user => user.animals)
  @JoinColumn({ name: 'user_id' })
  user: User;

  @ManyToOne(() => Lot, lot => lot.animals, { nullable: true })
  @JoinColumn({ name: 'lot_id' })
  lot?: Lot;

  @ManyToOne(() => Animal, animal => animal.children, { nullable: true })
  @JoinColumn({ name: 'mother_id' })
  mother?: Animal;

  @ManyToOne(() => Animal, animal => animal.children, { nullable: true })
  @JoinColumn({ name: 'father_id' })
  father?: Animal;

  @OneToMany(() => Animal, animal => animal.mother)
  children: Animal[];

  @OneToMany(() => AnimalVaccine, animalVaccine => animalVaccine.animal)
  vaccinations: AnimalVaccine[];

  @OneToMany(() => WeightRecord, weightRecord => weightRecord.animal)
  weightRecords: WeightRecord[];

  @OneToMany(() => PregnancyRecord, pregnancyRecord => pregnancyRecord.animal)
  pregnancies: PregnancyRecord[];

  @OneToMany(() => PregnancyRecord, pregnancyRecord => pregnancyRecord.father)
  offspring: PregnancyRecord[];

  @OneToMany(() => Sale, sale => sale.animal)
  sales: Sale[];

  // Virtual properties
  get age(): number {
    const today = new Date();
    const birthDate = new Date(this.birthDate);
    let age = today.getFullYear() - birthDate.getFullYear();
    const monthDiff = today.getMonth() - birthDate.getMonth();
    
    if (monthDiff < 0 || (monthDiff === 0 && today.getDate() < birthDate.getDate())) {
      age--;
    }
    
    return age;
  }

  get isPregnant(): boolean {
    return this.pregnancies?.some(p => p.pregnancyStatus === 'confirmed') || false;
  }
}