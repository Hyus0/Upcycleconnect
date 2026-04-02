import { request } from "./http";
import {
  getDashboardSnapshot,
  readCollection,
} from "./mockDb";

const GO_API_BASE = import.meta.env.VITE_GO_API_BASE ?? "/api-go";

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

async function readUsersFromApi() {
  const response = await request(`${GO_API_BASE}/admin/users`);
  return (response.items ?? []).map(normalizeUser);
}

async function readPrestationsFromApi(filters = {}) {
  const response = await request(`${GO_API_BASE}/admin/prestations${buildQuery({ type: filters.type })}`);
  return (response.items ?? []).map(normalizePrestation);
}

async function readCategoriesFromApi() {
  const response = await request(`${GO_API_BASE}/admin/categories`);
  return (response.items ?? []).map(normalizeCategory);
}

async function readEventsFromApi() {
  const response = await request(`${GO_API_BASE}/admin/events`);
  return (response.items ?? []).map(normalizeEvent);
}

async function readUsersLocal() {
  return (await readCollection("users")).map(normalizeUser);
}

async function readPrestationsLocal() {
  return (await readCollection("prestations")).map(normalizePrestation);
}

async function readCategoriesLocal() {
  return (await readCollection("categories")).map(normalizeCategory);
}

async function readEventsLocal() {
  return (await readCollection("events")).map(normalizeEvent);
}

export const capabilities = {
  dashboard: { metrics: true, localFallback: true },
  users: { list: true, detail: true, create: true, update: true, toggle: true, delete: true },
  prestations: { list: true, create: true, update: true, delete: true },
  categories: { list: true, create: true, update: true, delete: true },
  events: { list: true, create: true, update: true, delete: true }
};

export const adminApi = {
  async getDashboard() {
    try {
      const [metricsResponse, users, prestations] = await Promise.all([
        request(`${GO_API_BASE}/admin/metrics`),
        readUsersFromApi(),
        readPrestationsFromApi()
      ]);

      return buildDashboard({
        users,
        prestations,
        categories: await readCategoriesFromApi(),
        events: await readEventsFromApi(),
        metrics: metricsResponse.metrics ?? {},
        source: "api"
      });
    } catch {
      const db = await getDashboardSnapshot();
      return buildDashboard({
        users: (db.users ?? []).map(normalizeUser),
        prestations: (db.prestations ?? []).map(normalizePrestation),
        categories: (db.categories ?? []).map(normalizeCategory),
        events: (db.events ?? []).map(normalizeEvent),
        metrics: {},
        source: "local"
      });
    }
  },

  async listUsers(filters = {}) {
    let items;
    try {
      items = await readUsersFromApi();
    } catch {
      items = await readUsersLocal();
    }
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

  async listPrestations(filters = {}) {
    let items;
    try {
      items = await readPrestationsFromApi(filters);
    } catch {
      items = await readPrestationsLocal();
    }
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

  async listCategories(filters = {}) {
    let items;
    try {
      items = await readCategoriesFromApi();
    } catch {
      items = await readCategoriesLocal();
    }
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

  async listEvents(filters = {}) {
    let items;
    try {
      items = await readEventsFromApi();
    } catch {
      items = await readEventsLocal();
    }
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
