import jsPDF from 'jspdf';
import 'jspdf-autotable';
import { Animal } from '../../pages/contents/AnimalTable/types/type';
import { Sale } from '../../types/sale';
import { MilkCollection } from '../../types/milk-collection';
import { Reproduction } from '../../types/reproduction';

interface AutoTableOptions {
  startY: number;
  head: string[][];
  body: string[][];
  styles: {
    fontSize: number;
    cellPadding: number;
  };
  headStyles: {
    fillColor: number[];
    textColor: number[];
    fontStyle: string;
  };
  alternateRowStyles: {
    fillColor: number[];
  };
  margin: {
    left: number;
    right: number;
  };
}

interface JSPDFWithAutoTable extends jsPDF {
  autoTable: (options: AutoTableOptions) => void;
  lastAutoTable: {
    finalY: number;
  };
}

interface Translations {
  animalHistoryExport: {
    title: string;
    animalInfo: string;
    fields: Record<string, string>;
    statistics: string;
    stats: Record<string, string>;
    salesHistory: string;
    salesTable: Record<string, string>;
    milkHistory: string;
    milkTable: Record<string, string>;
    reproductionHistory: string;
    reproductionTable: Record<string, string>;
    footer: {
      page: string;
      of: string;
    };
    sex?: {
      male?: string;
      female?: string;
    };
    status?: {
      active?: string;
      inactive?: string;
      sold?: string;
      deceased?: string;
    };
  };
  common?: {
    yes?: string;
    no?: string;
    notInformed?: string;
  };
}

interface PDFGeneratorOptions {
  animal: Animal;
  sales: Sale[];
  milkCollections: MilkCollection[];
  reproductions: Reproduction[];
  animalImage?: string;
  translations: Translations;
}

export class AnimalHistoryPDFGenerator {
  private doc: JSPDFWithAutoTable;
  private pageWidth: number;
  private margin: number;
  private currentY: number;
  private translations: Translations;

  constructor(options: PDFGeneratorOptions) {
    this.doc = new jsPDF() as JSPDFWithAutoTable;
    this.pageWidth = this.doc.internal.pageSize.getWidth();
    this.margin = 20;
    this.currentY = 30;
    this.translations = options.translations;
    
    this.generatePDF(options);
  }

  private generatePDF(options: PDFGeneratorOptions) {
    const { animal, sales, milkCollections, reproductions, animalImage } = options;
    
    this.addHeader();
    this.addAnimalImage(animalImage);
    this.addAnimalInfo(animal);
    this.addStatistics(sales, milkCollections, reproductions);
    this.addSalesHistory(sales);
    this.addMilkHistory(milkCollections);
    this.addReproductionHistory(reproductions);
    this.addFooter();
    
    this.savePDF(animal);
  }

  private addHeader() {
    this.doc.setFontSize(24);
    this.doc.setFont('helvetica', 'bold');
    this.doc.setTextColor(66, 139, 202);
    this.doc.text(this.translations.animalHistoryExport.title, this.pageWidth / 2, this.currentY, { align: 'center' });
    
    this.currentY += 20;
    
    this.doc.setDrawColor(66, 139, 202);
    this.doc.setLineWidth(2);
    this.doc.line(this.margin, this.currentY, this.pageWidth - this.margin, this.currentY);
    
    this.currentY += 15;
  }

  private addAnimalImage(imageUrl?: string) {
    if (imageUrl) {
      try {
        const imgWidth = 50;
        const imgHeight = 50;
        const imgX = this.pageWidth - this.margin - imgWidth;
        const imgY = this.currentY;
        
        const imageFormat = imageUrl.includes('data:image/') ? 
          imageUrl.split('data:image/')[1].split(';')[0] : 'JPEG';
        
        this.doc.addImage(imageUrl, imageFormat.toUpperCase(), imgX, imgY, imgWidth, imgHeight);
        
        this.doc.setDrawColor(200, 200, 200);
        this.doc.setLineWidth(1);
        this.doc.rect(imgX, imgY, imgWidth, imgHeight);
        
        this.currentY += imgHeight + 10;
      } catch (error) {
        console.warn('Erro ao adicionar imagem do animal:', error);
      }
    }
  }

  private addAnimalInfo(animal: Animal) {
    this.doc.setFontSize(16);
    this.doc.setFont('helvetica', 'bold');
    this.doc.setTextColor(0, 0, 0);
    this.doc.text(this.translations.animalHistoryExport.animalInfo, this.margin, this.currentY);
    this.currentY += 15;

    const basicInfo = [
      [this.translations.animalHistoryExport.fields.name, animal.animal_name],
      [this.translations.animalHistoryExport.fields.earTag, animal.ear_tag_number_local.toString()],
      [this.translations.animalHistoryExport.fields.breed, animal.breed],
      [this.translations.animalHistoryExport.fields.type, animal.type],
      [this.translations.animalHistoryExport.fields.birthDate, animal.birth_date ? new Date(animal.birth_date).toLocaleDateString('pt-BR') : (this.translations.common?.notInformed || 'Não informado')],
    ];

    this.doc.setFont('helvetica', 'normal');
    this.doc.setFontSize(12);
    
    basicInfo.forEach(([label, value]) => {
      this.doc.setFont('helvetica', 'bold');
      this.doc.text(`${label}:`, this.margin, this.currentY);
      this.doc.setFont('helvetica', 'normal');
      this.doc.text(value, this.margin + 60, this.currentY);
      this.currentY += 8;
    });

    this.currentY += 10;

    const additionalInfo = [
      [this.translations.animalHistoryExport.fields.sex, this.getSexText(animal.sex)],
      [this.translations.animalHistoryExport.fields.status, this.getStatusText(animal.status)],
      [this.translations.animalHistoryExport.fields.confinement, animal.confinement ? this.translations.common?.yes || 'Sim' : this.translations.common?.no || 'Não'],
      [this.translations.animalHistoryExport.fields.fertilization, animal.fertilization ? this.translations.common?.yes || 'Sim' : this.translations.common?.no || 'Não'],
      [this.translations.animalHistoryExport.fields.castrated, animal.castrated ? this.translations.common?.yes || 'Sim' : this.translations.common?.no || 'Não'],
    ];

    if (animal.current_weight) {
      additionalInfo.push([this.translations.animalHistoryExport.fields.currentWeight, `${animal.current_weight} kg`]);
    }
    if (animal.ideal_weight) {
      additionalInfo.push([this.translations.animalHistoryExport.fields.idealWeight, `${animal.ideal_weight} kg`]);
    }
    if (animal.milk_production) {
      additionalInfo.push([this.translations.animalHistoryExport.fields.milkProduction, `${animal.milk_production} L/dia`]);
    }

    additionalInfo.forEach(([label, value]) => {
      this.doc.setFont('helvetica', 'bold');
      this.doc.text(`${label}:`, this.margin, this.currentY);
      this.doc.setFont('helvetica', 'normal');
      this.doc.text(value, this.margin + 60, this.currentY);
      this.currentY += 8;
    });

    this.currentY += 15;
  }

  private addStatistics(sales: Sale[], milkCollections: MilkCollection[], reproductions: Reproduction[]) {
    this.doc.setFontSize(14);
    this.doc.setFont('helvetica', 'bold');
    this.doc.setTextColor(66, 139, 202);
    this.doc.text(this.translations.animalHistoryExport.statistics, this.margin, this.currentY);
    this.currentY += 10;

    const totalSales = sales.length;
    const totalSalesValue = sales.reduce((sum, sale) => sum + sale.price, 0);
    const totalMilkCollections = milkCollections.length;
    const totalMilkQuantity = milkCollections.reduce((sum, collection) => sum + collection.quantity, 0);
    const totalReproductions = reproductions.length;

    const stats = [
      [this.translations.animalHistoryExport.stats.totalSales, totalSales.toString()],
      [this.translations.animalHistoryExport.stats.totalSalesValue, `R$ ${totalSalesValue.toFixed(2)}`],
      [this.translations.animalHistoryExport.stats.totalMilkCollections, totalMilkCollections.toString()],
      [this.translations.animalHistoryExport.stats.totalMilkQuantity, `${totalMilkQuantity.toFixed(2)} L`],
      [this.translations.animalHistoryExport.stats.totalReproductions, totalReproductions.toString()],
    ];

    this.doc.setFont('helvetica', 'normal');
    this.doc.setFontSize(11);
    this.doc.setTextColor(0, 0, 0);

    stats.forEach(([label, value]) => {
      this.doc.setFont('helvetica', 'bold');
      this.doc.text(`${label}:`, this.margin, this.currentY);
      this.doc.setFont('helvetica', 'normal');
      this.doc.text(value, this.margin + 80, this.currentY);
      this.currentY += 7;
    });

    this.currentY += 15;
  }

  private addSalesHistory(sales: Sale[]) {
    if (sales.length === 0) return;

    this.checkPageBreak(30);
    
    this.doc.setFontSize(14);
    this.doc.setFont('helvetica', 'bold');
    this.doc.setTextColor(66, 139, 202);
    this.doc.text(this.translations.animalHistoryExport.salesHistory, this.margin, this.currentY);
    this.currentY += 10;

    const salesData = sales.map(sale => [
      new Date(sale.sale_date).toLocaleDateString('pt-BR'),
      sale.buyer_name,
      `R$ ${sale.price.toFixed(2)}`,
      sale.notes || '-'
    ]);

    this.doc.autoTable({
      startY: this.currentY,
      head: [[
        this.translations.animalHistoryExport.salesTable.date,
        this.translations.animalHistoryExport.salesTable.buyer,
        this.translations.animalHistoryExport.salesTable.price,
        this.translations.animalHistoryExport.salesTable.notes
      ]],
      body: salesData,
      styles: { 
        fontSize: 10,
        cellPadding: 4
      },
      headStyles: { 
        fillColor: [66, 139, 202],
        textColor: [255, 255, 255],
        fontStyle: 'bold'
      },
      alternateRowStyles: {
        fillColor: [245, 245, 245]
      },
      margin: { left: this.margin, right: this.margin }
    });

    this.currentY = this.doc.lastAutoTable.finalY + 15;
  }

  private addMilkHistory(milkCollections: MilkCollection[]) {
    if (milkCollections.length === 0) return;

    this.checkPageBreak(30);
    
    this.doc.setFontSize(14);
    this.doc.setFont('helvetica', 'bold');
    this.doc.setTextColor(34, 139, 34);
    this.doc.text(this.translations.animalHistoryExport.milkHistory, this.margin, this.currentY);
    this.currentY += 10;

    const milkData = milkCollections.map((collection: MilkCollection) => [
      new Date(collection.collection_date).toLocaleDateString('pt-BR'),
      `${collection.quantity}L`,
      collection.quality || '-',
      collection.notes || '-'
    ]);

    this.doc.autoTable({
      startY: this.currentY,
      head: [[
        this.translations.animalHistoryExport.milkTable.date,
        this.translations.animalHistoryExport.milkTable.quantity,
        this.translations.animalHistoryExport.milkTable.quality,
        this.translations.animalHistoryExport.milkTable.notes
      ]],
      body: milkData,
      styles: { 
        fontSize: 10,
        cellPadding: 4
      },
      headStyles: { 
        fillColor: [34, 139, 34],
        textColor: [255, 255, 255],
        fontStyle: 'bold'
      },
      alternateRowStyles: {
        fillColor: [245, 255, 245]
      },
      margin: { left: this.margin, right: this.margin }
    });

    this.currentY = this.doc.lastAutoTable.finalY + 15;
  }

  private addReproductionHistory(reproductions: Reproduction[]) {
    if (reproductions.length === 0) return;

    this.checkPageBreak(30);
    
    this.doc.setFontSize(14);
    this.doc.setFont('helvetica', 'bold');
    this.doc.setTextColor(255, 140, 0);
    this.doc.text(this.translations.animalHistoryExport.reproductionHistory, this.margin, this.currentY);
    this.currentY += 10;

    const reproductionData = reproductions.map((reproduction: Reproduction) => [
      new Date(reproduction.date).toLocaleDateString('pt-BR'),
      reproduction.phase,
      reproduction.notes || '-'
    ]);

    this.doc.autoTable({
      startY: this.currentY,
      head: [[
        this.translations.animalHistoryExport.reproductionTable.date,
        this.translations.animalHistoryExport.reproductionTable.phase,
        this.translations.animalHistoryExport.reproductionTable.notes
      ]],
      body: reproductionData,
      styles: { 
        fontSize: 10,
        cellPadding: 4
      },
      headStyles: { 
        fillColor: [255, 140, 0],
        textColor: [255, 255, 255],
        fontStyle: 'bold'
      },
      alternateRowStyles: {
        fillColor: [255, 248, 240]
      },
      margin: { left: this.margin, right: this.margin }
    });

    this.currentY = this.doc.lastAutoTable.finalY + 15;
  }

  private addFooter() {
    const pageCount = this.doc.getNumberOfPages();
    
    for (let i = 1; i <= pageCount; i++) {
      this.doc.setPage(i);
      this.doc.setFontSize(8);
      this.doc.setFont('helvetica', 'normal');
      this.doc.setTextColor(128, 128, 128);
      
      this.doc.text(
        `${this.translations.animalHistoryExport.footer.page} ${i} ${this.translations.animalHistoryExport.footer.of} ${pageCount}`,
        this.pageWidth / 2,
        this.doc.internal.pageSize.getHeight() - 10,
        { align: 'center' }
      );
      
      this.doc.text(
        new Date().toLocaleDateString('pt-BR'),
        this.pageWidth - this.margin,
        this.doc.internal.pageSize.getHeight() - 10,
        { align: 'right' }
      );
    }
  }

  private checkPageBreak(requiredSpace: number) {
    const pageHeight = this.doc.internal.pageSize.getHeight();
    if (this.currentY + requiredSpace > pageHeight - 30) {
      this.doc.addPage();
      this.currentY = 30;
    }
  }

  private savePDF(animal: Animal) {
    const fileName = `${animal.animal_name}_historico_completo_${new Date().toISOString().split('T')[0]}.pdf`;
    this.doc.save(fileName);
  }

  private getSexText(sex: number): string {
    return sex === 0 ? (this.translations.animalHistoryExport.sex?.male || 'Macho') : (this.translations.animalHistoryExport.sex?.female || 'Fêmea');
  }

  private getStatusText(status: number): string {
    const statusMap: { [key: number]: string } = {
      0: this.translations.animalHistoryExport.status?.active || 'Ativo',
      1: this.translations.animalHistoryExport.status?.inactive || 'Inativo',
      2: this.translations.animalHistoryExport.status?.sold || 'Vendido',
      3: this.translations.animalHistoryExport.status?.deceased || 'Falecido'
    };
    return statusMap[status] || (this.translations.common?.notInformed || 'Não informado');
  }
}

export const generateAnimalHistoryPDF = (options: PDFGeneratorOptions) => {
  return new AnimalHistoryPDFGenerator(options);
};