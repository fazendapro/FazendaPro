export const translations = {
  sales: {
    title: "Sales",
    createSale: "Register Sale",
    history: "Sales History",
    summary: "Summary",
    summaryDetails: {
      totalSales: "Total Sales",
      totalValue: "Total Value"
    },
    table: {
      animal: "Animal",
      earTag: "Ear Tag",
      status: "Status",
      buyer: "Buyer",
      price: "Price",
      saleDate: "Sale Date",
      notes: "Notes",
      items: "items"
    },
    filters: {
      startDate: "Start Date",
      endDate: "End Date",
      apply: "Apply Filters",
      clear: "Clear Filters"
    },
    export: {
      title: "Export",
      comingSoon: "Export functionality coming soon"
    },
    confirmDelete: {
      title: "Confirm Deletion",
      content: "Are you sure you want to delete the sale to buyer {{buyer}}?"
    }
  },
  saleModal: {
    title: {
      create: "Register Sale",
      edit: "Edit Sale"
    },
    fields: {
      animal: "Animal",
      buyerName: "Buyer Name",
      price: "Price",
      saleDate: "Sale Date",
      notes: "Notes"
    },
    placeholders: {
      animal: "Select an animal",
      buyerName: "Enter buyer name",
      price: "Enter price",
      saleDate: "Select sale date",
      notes: "Enter notes (optional)"
    },
    validation: {
      animalRequired: "Animal selection is required",
      buyerNameRequired: "Buyer name is required",
      priceRequired: "Price is required",
      priceMin: "Price must be greater than zero",
      saleDateRequired: "Sale date is required"
    },
    success: {
      created: "Sale registered successfully!",
      updated: "Sale updated successfully!"
    }
  },
  animalHistoryExport: {
    title: "Complete Animal History",
    animalInfo: "Animal Information",
    fields: {
      name: "Name",
      earTag: "Ear Tag",
      breed: "Breed",
      type: "Type",
      birthDate: "Birth Date",
      sex: "Sex",
      status: "Status",
      confinement: "Confinement",
      fertilization: "Fertilization",
      castrated: "Castrated",
      currentWeight: "Current Weight",
      idealWeight: "Ideal Weight",
      milkProduction: "Milk Production"
    },
    sex: {
      male: "Male",
      female: "Female"
    },
    status: {
      active: "Active",
      inactive: "Inactive",
      sold: "Sold",
      deceased: "Deceased"
    },
    statistics: "Statistics",
    stats: {
      totalSales: "Total Sales",
      totalSalesValue: "Total Sales Value",
      totalMilkCollections: "Total Milk Collections",
      totalMilkQuantity: "Total Milk Quantity",
      totalReproductions: "Total Reproductions"
    },
    salesHistory: "Sales History",
    salesTable: {
      date: "Date",
      buyer: "Buyer",
      price: "Price",
      notes: "Notes"
    },
    milkHistory: "Milk Collection History",
    milkTable: {
      date: "Date",
      quantity: "Quantity",
      quality: "Quality",
      notes: "Notes"
    },
    reproductionHistory: "Reproduction History",
    reproductionTable: {
      date: "Date",
      phase: "Phase",
      notes: "Notes"
    },
    footer: {
      page: "Page",
      of: "of"
    },
    units: {
      kg: "kg",
      liters: "L",
      litersPerDay: "L/day",
      currency: "$"
    },
    dateFormat: "en-US",
    fileName: "complete_history",
    success: "PDF generated successfully!",
    error: "Error generating PDF"
  },
  common: {
    actions: "Actions",
    edit: "Edit",
    delete: "Delete",
    yes: "Yes",
    no: "No",
    cancel: "Cancel",
    save: "Save",
    create: "Create",
    update: "Update",
    notInformed: "Not informed"
  }
};

