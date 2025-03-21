import { MigrationInterface, QueryRunner } from 'typeorm';

export class Init1612345678901 implements MigrationInterface {
  public async up(queryRunner: QueryRunner): Promise<void> {
    await queryRunner.query(`
      CREATE TABLE \`user\` (
        \`id\` INT AUTO_INCREMENT PRIMARY KEY,
        \`firstName\` VARCHAR(255) NOT NULL,
        \`lastName\` VARCHAR(255) NOT NULL,
        \`email\` VARCHAR(255) NOT NULL,
        \`phone\` VARCHAR(255) NOT NULL,
        \`password\` VARCHAR(255) NOT NULL,
        \`lastAccess\` DATETIME NOT NULL,
        \`farmId\` INT NOT NULL,
        \`isActive\` TINYINT DEFAULT 1
      )
    `);
  }

  public async down(queryRunner: QueryRunner): Promise<void> {
    await queryRunner.query(`DROP TABLE \`user\``);
  }
}
