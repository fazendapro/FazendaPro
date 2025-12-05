export const translations = {
  sales: {
    title: "Vendas",
    createSale: "Registrar Venda",
    history: "Histórico de Vendas",
    summary: "Resumo",
    summaryDetails: {
      totalSales: "Total de Vendas",
      totalValue: "Valor Total"
    },
    table: {
      animal: "Animal",
      earTag: "Brinco",
      status: "Status",
      buyer: "Comprador",
      price: "Preço",
      saleDate: "Data da Venda",
      notes: "Observações",
      items: "itens"
    },
    filters: {
      startDate: "Data Inicial",
      endDate: "Data Final",
      apply: "Aplicar Filtros",
      clear: "Limpar Filtros"
    },
    export: {
      title: "Exportar",
      comingSoon: "Funcionalidade de exportação em breve"
    },
    confirmDelete: {
      title: "Confirmar Exclusão",
      content: "Tem certeza que deseja excluir a venda do comprador {{buyer}}?"
    }
  },
  saleModal: {
    title: {
      create: "Registrar Venda",
      edit: "Editar Venda"
    },
    fields: {
      animal: "Animal",
      buyerName: "Nome do Comprador",
      price: "Preço",
      saleDate: "Data da Venda",
      notes: "Observações"
    },
        placeholders: {
          animal: "Selecione um animal",
          buyerName: "Digite o nome do comprador",
          price: "Digite o preço",
          saleDate: "Selecione a data da venda",
          notes: "Digite observações (opcional)"
        },
    validation: {
      animalRequired: "Seleção do animal é obrigatória",
      buyerNameRequired: "Nome do comprador é obrigatório",
      priceRequired: "Preço é obrigatório",
      priceMin: "Preço deve ser maior que zero",
      saleDateRequired: "Data da venda é obrigatória"
    },
    success: {
      created: "Venda registrada com sucesso!",
      updated: "Venda atualizada com sucesso!"
    }
  },
  animalHistoryExport: {
    title: "Histórico Completo do Animal",
    animalInfo: "Informações do Animal",
    fields: {
      name: "Nome",
      earTag: "Brinco",
      breed: "Raça",
      type: "Tipo",
      birthDate: "Data de Nascimento",
      sex: "Sexo",
      status: "Status",
      confinement: "Confinamento",
      fertilization: "Fertilização",
      castrated: "Castrado",
      currentWeight: "Peso Atual",
      idealWeight: "Peso Ideal",
      milkProduction: "Produção de Leite"
    },
    sex: {
      male: "Macho",
      female: "Fêmea"
    },
    status: {
      active: "Ativo",
      inactive: "Inativo",
      sold: "Vendido",
      deceased: "Falecido"
    },
    statistics: "Estatísticas",
    stats: {
      totalSales: "Total de Vendas",
      totalSalesValue: "Valor Total das Vendas",
      totalMilkCollections: "Total de Ordenhas",
      totalMilkQuantity: "Quantidade Total de Leite",
      totalReproductions: "Total de Reproduções"
    },
    salesHistory: "Histórico de Vendas",
    salesTable: {
      date: "Data",
      buyer: "Comprador",
      price: "Preço",
      notes: "Observações"
    },
    milkHistory: "Histórico de Ordenha",
    milkTable: {
      date: "Data",
      quantity: "Quantidade",
      quality: "Qualidade",
      notes: "Observações"
    },
    reproductionHistory: "Histórico de Reprodução",
    reproductionTable: {
      date: "Data",
      phase: "Fase",
      notes: "Observações"
    },
    footer: {
      page: "Página",
      of: "de"
    },
    units: {
      kg: "kg",
      liters: "L",
      litersPerDay: "L/dia",
      currency: "R$"
    },
    dateFormat: "pt-BR",
    fileName: "historico_completo",
    success: "PDF gerado com sucesso!",
    error: "Erro ao gerar PDF"
  },
  common: {
    actions: "Ações",
    edit: "Editar",
    delete: "Excluir",
    yes: "Sim",
    no: "Não",
    cancel: "Cancelar",
    save: "Salvar",
    create: "Criar",
    update: "Atualizar",
    notInformed: "Não informado"
  }
};
