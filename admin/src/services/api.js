import { request } from "./http";

const GO_API_BASE = import.meta.env.VITE_GO_API_BASE ?? "http://localhost:8081/api";
const GO_PUBLIC_BASE = import.meta.env.VITE_GO_PUBLIC_BASE ?? "http://localhost:8081";

function buildQuery(params = {}) {
  const searchParams = new URLSearchParams();
  Object.entries(params).forEach(([key, value]) => {
    if (value !== undefined && value !== null && value !== "") {
      searchParams.set(key, String(value));
    }
  });
  const query = searchParams.toString();
  return query ? `?${query}` : "";
}

function paginate(items, page = 1, pageSize = 10) {
  const currentPage = Number(page) || 1;
  const size = Number(pageSize) || 10;
  const start = (currentPage - 1) * size;
  return {
    items: items.slice(start, start + size),
    pagination: {
      page: currentPage,
      pageSize: size,
      total: items.length,
      totalPages: Math.max(1, Math.ceil(items.length / size))
    }
  };
}

function normalizeUser(row) {
  const firstName = row.prenom ?? row.firstName ?? "";
  const lastName = row.nom ?? row.lastName ?? "";
  return {
    id: row.id,
    firstName,
    lastName,
    fullName: row.name ?? `${firstName} ${lastName}`.trim(),
    email: row.mail ?? row.email ?? "",
    address: row.adresse ?? row.address ?? "",
    city: row.ville ?? row.city ?? "",
    postalCode: row.code_postal ?? row.postalCode ?? "",
    birthDate: row.date_naissance ?? row.birthDate ?? "",
    createdAt: row.date_inscription ?? row.createdAt ?? "",
    role: (row.role ?? "Particulier").toString(),
    status: row.status ?? "active",
    languageId: Number(row.id_langue ?? row.languageId ?? 1)
  };
}

function normalizePrestation(row) {
  return {
    id: row.id,
    title: row.title ?? row.titre ?? "",
    description: row.description ?? "",
    type: row.type ?? "service",
    price: Number(row.price ?? row.prix ?? 0),
    status: row.status ?? row.statut ?? "draft",
    provider: row.provider ?? "Non assigne",
    createdAt: row.createdAt ?? row.date_creation ?? ""
  };
}

function normalizeCategory(row) {
  return {
    id: row.id,
    name: row.name ?? "",
    parentId: row.parentId ?? "",
    description: row.description ?? "",
    status: row.status ?? "active"
  };
}

function normalizeEvent(row) {
  return {
    id: row.id,
    title: row.title ?? "",
    location: row.location ?? "",
    date: row.date ?? "",
    status: row.status ?? "planned",
    capacity: Number(row.capacity ?? 0),
    description: row.description ?? ""
  };
}

function normalizeAnnonce(row) {
  return {
    id: row.id,
    sellerId: Number(row.id_vendeur ?? row.sellerId ?? 1),
    buyerId: row.id_acheteur ?? row.buyerId ?? null,
    title: row.titre ?? row.title ?? "",
    description: row.description ?? "",
    status: row.statut ?? row.status ?? "Disponible",
    validation: row.est_valide ?? row.validation ?? "En attente",
    price: Number(row.prix ?? row.price ?? 0),
    condition: row.etat_objet ?? row.condition ?? "",
    address: row.adresse ?? row.address ?? "",
    city: row.ville ?? row.city ?? "",
    postalCode: row.code_postal ?? row.postalCode ?? "",
    createdAt: row.date_creation ?? row.createdAt ?? "",
    type: (row.type ?? "Don").toString()
  };
}

function normalizeFinanceRecord(row) {
  return {
    id: row.id,
    label: row.label ?? "",
    category: row.category ?? "",
    amount: Number(row.amount ?? 0),
    status: row.status ?? "pending",
    dueDate: row.dueDate ?? "",
    source: row.source ?? ""
  };
}

function normalizeNotification(row) {
  return {
    id: row.id,
    title: row.title ?? "",
    channel: row.channel ?? "email",
    audience: row.audience ?? "all",
    status: row.status ?? "draft",
    scheduledAt: row.scheduledAt ?? "",
    message: row.message ?? ""
  };
}

function filterUsers(items, filters = {}) {
  let filtered = items;
  if (filters.role) {
    filtered = filtered.filter((user) => user.role.toLowerCase() === filters.role.toLowerCase());
  }
  if (filters.status) {
    filtered = filtered.filter((user) => user.status === filters.status);
  }
  if (filters.search) {
    const needle = filters.search.toLowerCase();
    filtered = filtered.filter((user) =>
      [user.fullName, user.email, user.city].join(" ").toLowerCase().includes(needle)
    );
  }
  return filtered;
}

function filterPrestations(items, filters = {}) {
  let filtered = items;
  if (filters.type) {
    filtered = filtered.filter((item) => item.type === filters.type);
  }
  if (filters.search) {
    const needle = filters.search.toLowerCase();
    filtered = filtered.filter((item) =>
      [item.title, item.description, item.provider].join(" ").toLowerCase().includes(needle)
    );
  }
  if (filters.status) {
    filtered = filtered.filter((item) => item.status === filters.status);
  }
  if (filters.sortBy === "price-desc") {
    filtered = filtered.slice().sort((a, b) => b.price - a.price);
  } else if (filters.sortBy === "price-asc") {
    filtered = filtered.slice().sort((a, b) => a.price - b.price);
  } else {
    filtered = filtered.slice().sort((a, b) => a.title.localeCompare(b.title));
  }
  return filtered;
}

function filterCategories(items, filters = {}) {
  let filtered = items;
  if (filters.search) {
    const needle = filters.search.toLowerCase();
    filtered = filtered.filter((item) =>
      [item.name, item.description].join(" ").toLowerCase().includes(needle)
    );
  }
  return filtered;
}

function filterEvents(items, filters = {}) {
  let filtered = items;
  if (filters.search) {
    const needle = filters.search.toLowerCase();
    filtered = filtered.filter((item) =>
      [item.title, item.location, item.description].join(" ").toLowerCase().includes(needle)
    );
  }
  if (filters.status) {
    filtered = filtered.filter((item) => item.status === filters.status);
  }
  if (filters.date) {
    filtered = filtered.filter((item) => item.date === filters.date);
  }
  return filtered.sort((a, b) => a.date.localeCompare(b.date));
}

function filterAnnonces(items, filters = {}) {
  let filtered = items;
  if (filters.type) {
    filtered = filtered.filter((item) => item.type.toLowerCase() === filters.type.toLowerCase());
  }
  if (filters.validation) {
    filtered = filtered.filter((item) => item.validation.toLowerCase() === filters.validation.toLowerCase());
  }
  if (filters.search) {
    const needle = filters.search.toLowerCase();
    filtered = filtered.filter((item) =>
      [item.title, item.description, item.city, item.postalCode].join(" ").toLowerCase().includes(needle)
    );
  }
  return filtered.slice().sort((a, b) => Number(b.id) - Number(a.id));
}

function filterModeration(items, filters = {}) {
  let filtered = items;
  if (filters.type) {
    filtered = filtered.filter((item) => item.type === filters.type);
  }
  if (filters.status) {
    filtered = filtered.filter((item) => item.status === filters.status);
  }
  if (filters.search) {
    const needle = filters.search.toLowerCase();
    filtered = filtered.filter((item) =>
      [item.title, item.owner, item.description].join(" ").toLowerCase().includes(needle)
    );
  }
  return filtered;
}

function filterFinance(items, filters = {}) {
  let filtered = items;
  if (filters.status) {
    filtered = filtered.filter((item) => item.status === filters.status);
  }
  if (filters.category) {
    filtered = filtered.filter((item) => item.category === filters.category);
  }
  return filtered.sort((a, b) => (b.dueDate ?? "").localeCompare(a.dueDate ?? ""));
}

function filterNotifications(items, filters = {}) {
  let filtered = items;
  if (filters.status) {
    filtered = filtered.filter((item) => item.status === filters.status);
  }
  if (filters.channel) {
    filtered = filtered.filter((item) => item.channel === filters.channel);
  }
  if (filters.search) {
    const needle = filters.search.toLowerCase();
    filtered = filtered.filter((item) =>
      [item.title, item.message, item.audience].join(" ").toLowerCase().includes(needle)
    );
  }
  return filtered;
}

async function readUsersFromApi() {
  const response = await request(`${GO_API_BASE}/admin/users`);
  return (response.items ?? []).map(normalizeUser);
}

async function readPrestationsFromApi(filters = {}) {
  const response = await request(`${GO_API_BASE}/admin/prestations${buildQuery({ type: filters.type })}`);
  return (response.items ?? []).map(normalizePrestation);
}

async function readAnnoncesFromApi() {
  const response = await request(`${GO_PUBLIC_BASE}/annonces`);
  return (Array.isArray(response) ? response : []).map(normalizeAnnonce);
}

async function readCategoriesFromApi() {
  const response = await request(`${GO_API_BASE}/admin/categories`);
  return (response.items ?? []).map(normalizeCategory);
}

async function readEventsFromApi() {
  const response = await request(`${GO_API_BASE}/admin/events`);
  return (response.items ?? []).map(normalizeEvent);
}

async function readModerationQueueFromApi() {
  const response = await request(`${GO_API_BASE}/admin/moderation/queue`);
  return response.items ?? [];
}

async function readFinanceFromApi() {
  const response = await request(`${GO_API_BASE}/admin/finance/overview`);
  return {
    summary: response.summary ?? {},
    items: (response.items ?? []).map(normalizeFinanceRecord)
  };
}

async function readNotificationsFromApi() {
  const response = await request(`${GO_API_BASE}/admin/notifications`);
  return (response.items ?? []).map(normalizeNotification);
}

export const capabilities = {
  dashboard: { metrics: true, localFallback: false },
  users: { list: true, detail: true, create: true, update: true, toggle: true, delete: true },
  prestations: { list: true, detail: true, create: true, update: true, delete: true },
  annonces: { list: true, create: true },
  categories: { list: true, detail: true, create: true, update: true, delete: true },
  events: { list: true, detail: true, create: true, update: true, delete: true },
  moderation: { queue: true, publish: true, archive: true },
  finance: { overview: true },
  notifications: { list: true, create: true, update: true, delete: true }
};

export const adminApi = {
  async getDashboard() {
    const [metricsResponse, users, prestations, categories, events] = await Promise.all([
      request(`${GO_API_BASE}/admin/metrics`),
      readUsersFromApi(),
      readPrestationsFromApi(),
      readCategoriesFromApi(),
      readEventsFromApi()
    ]);

    return buildDashboard({
      users,
      prestations,
      categories,
      events,
      metrics: metricsResponse.metrics ?? {},
      source: metricsResponse.source ?? "api"
    });
  },

  async listUsers(filters = {}) {
    const items = await readUsersFromApi();
    return paginate(filterUsers(items, filters), filters.page, filters.pageSize);
  },

  async createUser(payload) {
    const response = await request(`${GO_API_BASE}/admin/users`, {
      method: "POST",
      body: JSON.stringify({
        ...payload,
        createdAt: payload.createdAt ?? new Date().toISOString().slice(0, 10)
      })
    });
    return response.created ?? response;
  },

  async updateUser(id, payload) {
    const response = await request(`${GO_API_BASE}/admin/users/${id}`, {
      method: "PUT",
      body: JSON.stringify(payload)
    });
    return response.updated ?? response;
  },

  async toggleUserStatus(id) {
    const response = await request(`${GO_API_BASE}/admin/users/${id}/status`, {
      method: "PATCH"
    });
    return response.updated ?? response;
  },

  async deleteUser(id) {
    return request(`${GO_API_BASE}/admin/users/${id}`, { method: "DELETE" });
  },

  async getUser(id) {
    return request(`${GO_API_BASE}/admin/users/${id}`);
  },

  async listPrestations(filters = {}) {
    const items = await readPrestationsFromApi(filters);
    return paginate(filterPrestations(items, filters), filters.page, filters.pageSize);
  },

  async createPrestation(payload) {
    const requestPayload = {
      title: payload.title,
      description: payload.description,
      type: payload.type,
      price: Number(payload.price ?? 0),
      status: payload.status ?? "draft",
      provider: payload.provider ?? "Equipe locale"
    };

    const response = await request(`${GO_API_BASE}/admin/prestations`, {
      method: "POST",
      body: JSON.stringify({
        ...requestPayload,
        createdAt: new Date().toISOString().slice(0, 10)
      })
    });
    return response.created ?? response;
  },

  async updatePrestation(id, payload) {
    const response = await request(`${GO_API_BASE}/admin/prestations/${id}`, {
      method: "PUT",
      body: JSON.stringify(payload)
    });
    return response.updated ?? response;
  },

  async deletePrestation(id) {
    return request(`${GO_API_BASE}/admin/prestations/${id}`, { method: "DELETE" });
  },

  async getPrestation(id) {
    return request(`${GO_API_BASE}/admin/prestations/${id}`);
  },

  async listAnnonces(filters = {}) {
    const items = await readAnnoncesFromApi();
    return paginate(filterAnnonces(items, filters), filters.page, filters.pageSize);
  },

  async createAnnonce(payload) {
    const requestPayload = {
      id_vendeur: Number(payload.sellerId ?? 1),
      id_acheteur: null,
      titre: payload.title,
      description: payload.description,
      statut: payload.status,
      est_valide: payload.validation,
      prix: Number(payload.price ?? 0),
      etat_objet: payload.condition,
      adresse: payload.address,
      ville: payload.city,
      code_postal: payload.postalCode,
      type: payload.type
    };

    const response = await request(`${GO_PUBLIC_BASE}/annonces`, {
      method: "POST",
      body: JSON.stringify(requestPayload)
    });
    return response.created ?? response;
  },

  async listCategories(filters = {}) {
    const items = await readCategoriesFromApi();
    return paginate(filterCategories(items, filters), filters.page, filters.pageSize);
  },

  async createCategory(payload) {
    const response = await request(`${GO_API_BASE}/admin/categories`, {
      method: "POST",
      body: JSON.stringify(payload)
    });
    return response.created ?? response;
  },

  async updateCategory(id, payload) {
    const response = await request(`${GO_API_BASE}/admin/categories/${id}`, {
      method: "PUT",
      body: JSON.stringify(payload)
    });
    return response.updated ?? response;
  },

  async deleteCategory(id) {
    return request(`${GO_API_BASE}/admin/categories/${id}`, { method: "DELETE" });
  },

  async getCategory(id) {
    return request(`${GO_API_BASE}/admin/categories/${id}`);
  },

  async listEvents(filters = {}) {
    const items = await readEventsFromApi();
    return paginate(filterEvents(items, filters), filters.page, filters.pageSize);
  },

  async createEvent(payload) {
    const response = await request(`${GO_API_BASE}/admin/events`, {
      method: "POST",
      body: JSON.stringify(payload)
    });
    return response.created ?? response;
  },

  async updateEvent(id, payload) {
    const response = await request(`${GO_API_BASE}/admin/events/${id}`, {
      method: "PUT",
      body: JSON.stringify(payload)
    });
    return response.updated ?? response;
  },

  async deleteEvent(id) {
    return request(`${GO_API_BASE}/admin/events/${id}`, { method: "DELETE" });
  },

  async getEvent(id) {
    return request(`${GO_API_BASE}/admin/events/${id}`);
  },

  async listModerationQueue(filters = {}) {
    const items = await readModerationQueueFromApi();
    return paginate(filterModeration(items, filters), filters.page, filters.pageSize);
  },

  async publishModerationItem(type, id) {
    const response = await request(`${GO_API_BASE}/admin/moderation/${type}/${id}/publish`, {
      method: "PATCH"
    });
    return response.updated ?? response;
  },

  async archiveModerationItem(type, id) {
    const response = await request(`${GO_API_BASE}/admin/moderation/${type}/${id}/archive`, {
      method: "PATCH"
    });
    return response.updated ?? response;
  },

  async getFinanceOverview(filters = {}) {
    const response = await readFinanceFromApi();
    return {
      summary: response.summary,
      ...paginate(filterFinance(response.items, filters), filters.page, filters.pageSize)
    };
  },

  async listNotifications(filters = {}) {
    const items = await readNotificationsFromApi();
    return paginate(filterNotifications(items, filters), filters.page, filters.pageSize);
  },

  async createNotification(payload) {
    const response = await request(`${GO_API_BASE}/admin/notifications`, {
      method: "POST",
      body: JSON.stringify(payload)
    });
    return response.created ?? response;
  },

  async updateNotificationStatus(id, status) {
    const response = await request(`${GO_API_BASE}/admin/notifications/${id}/status`, {
      method: "PATCH",
      body: JSON.stringify({ status })
    });
    return response.updated ?? response;
  },

  async deleteNotification(id) {
    return request(`${GO_API_BASE}/admin/notifications/${id}`, { method: "DELETE" });
  },

  getCapabilities() {
    return capabilities;
  }
};

function buildDashboard({ users, prestations, categories, events, metrics, source }) {
  const usersByRole = users.reduce((acc, user) => {
    const key = user.role.toLowerCase();
    acc[key] = (acc[key] ?? 0) + 1;
    return acc;
  }, {});

  return {
    source,
    stats: [
      { label: "Utilisateurs", value: metrics.users ?? users.length, tone: "green" },
      { label: "Prestations", value: metrics.annonces ?? prestations.length, tone: "teal" },
      { label: "Categories", value: categories.length, tone: "sand" },
      { label: "Evenements", value: events.length, tone: "amber" }
    ],
    charts: {
      usersByRole: [
        { label: "Particuliers", value: usersByRole.particulier ?? 0 },
        { label: "Prestataires", value: usersByRole.prestataire ?? 0 },
        { label: "Admins", value: usersByRole.admin ?? 0 }
      ],
      resources: [
        { label: "Prestations", value: prestations.length },
        { label: "Categories", value: categories.length },
        { label: "Evenements", value: events.length }
      ]
    },
    quickNotes: [
      { tone: "green", text: `${users.length} utilisateurs disponibles.` },
      { tone: "green", text: `${prestations.length} prestations visibles.` },
      { tone: source === "local" ? "amber" : "green", text: source === "local" ? "Mode local actif." : "Donnees API actives." }
    ],
    recentActivity: [...prestations, ...events]
      .slice(0, 5)
      .map((item) => ({
        id: item.id,
        title: item.title,
        subtitle: item.description || item.location || "Activite"
      }))
  };
}
