export const translations = {
  sales: {
    title: "Ventas",
    createSale: "Registrar Venta",
    history: "Historial de Ventas",
    summary: "Resumen",
    summaryDetails: {
      totalSales: "Total de Ventas",
      totalValue: "Valor Total"
    },
    table: {
      animal: "Animal",
      earTag: "Arete",
      status: "Estado",
      buyer: "Comprador",
      price: "Precio",
      saleDate: "Fecha de Venta",
      notes: "Observaciones",
      items: "elementos"
    },
    filters: {
      startDate: "Fecha Inicial",
      endDate: "Fecha Final",
      apply: "Aplicar Filtros",
      clear: "Limpiar Filtros"
    },
    export: {
      title: "Exportar",
      comingSoon: "Funcionalidad de exportación próximamente"
    },
    confirmDelete: {
      title: "Confirmar Eliminación",
      content: "¿Está seguro de que desea eliminar la venta al comprador {{buyer}}?"
    }
  },
  saleModal: {
    title: {
      create: "Registrar Venta",
      edit: "Editar Venta"
    },
    fields: {
      animal: "Animal",
      buyerName: "Nombre del Comprador",
      price: "Precio",
      saleDate: "Fecha de Venta",
      notes: "Observaciones"
    },
    placeholders: {
      animal: "Seleccione un animal",
      buyerName: "Ingrese el nombre del comprador",
      price: "Ingrese el precio",
      saleDate: "Seleccione la fecha de venta",
      notes: "Ingrese observaciones (opcional)"
    },
    validation: {
      animalRequired: "La selección del animal es obligatoria",
      buyerNameRequired: "El nombre del comprador es obligatorio",
      priceRequired: "El precio es obligatorio",
      priceMin: "El precio debe ser mayor que cero",
      saleDateRequired: "La fecha de venta es obligatoria"
    },
    success: {
      created: "¡Venta registrada con éxito!",
      updated: "¡Venta actualizada con éxito!"
    }
  },
  animalHistoryExport: {
    title: "Historial Completo del Animal",
    animalInfo: "Información del Animal",
    fields: {
      name: "Nombre",
      earTag: "Arete",
      breed: "Raza",
      type: "Tipo",
      birthDate: "Fecha de Nacimiento",
      sex: "Sexo",
      status: "Estado",
      confinement: "Confinamiento",
      fertilization: "Fertilización",
      castrated: "Castrado",
      currentWeight: "Peso Actual",
      idealWeight: "Peso Ideal",
      milkProduction: "Producción de Leche"
    },
    sex: {
      male: "Macho",
      female: "Hembra"
    },
    status: {
      active: "Activo",
      inactive: "Inactivo",
      sold: "Vendido",
      deceased: "Fallecido"
    },
    statistics: "Estadísticas",
    stats: {
      totalSales: "Total de Ventas",
      totalSalesValue: "Valor Total de las Ventas",
      totalMilkCollections: "Total de Ordeños",
      totalMilkQuantity: "Cantidad Total de Leche",
      totalReproductions: "Total de Reproducciones"
    },
    salesHistory: "Historial de Ventas",
    salesTable: {
      date: "Fecha",
      buyer: "Comprador",
      price: "Precio",
      notes: "Observaciones"
    },
    milkHistory: "Historial de Ordeño",
    milkTable: {
      date: "Fecha",
      quantity: "Cantidad",
      quality: "Calidad",
      notes: "Observaciones"
    },
    reproductionHistory: "Historial de Reproducción",
    reproductionTable: {
      date: "Fecha",
      phase: "Fase",
      notes: "Observaciones"
    },
    footer: {
      page: "Página",
      of: "de"
    },
    units: {
      kg: "kg",
      liters: "L",
      litersPerDay: "L/día",
      currency: "$"
    },
    dateFormat: "es-ES",
    fileName: "historial_completo",
    success: "¡PDF generado con éxito!",
    error: "Error al generar PDF"
  },
  common: {
    actions: "Acciones",
    edit: "Editar",
    delete: "Eliminar",
    yes: "Sí",
    no: "No",
    cancel: "Cancelar",
    save: "Guardar",
    create: "Crear",
    update: "Actualizar",
    notInformed: "No informado"
  }
};

