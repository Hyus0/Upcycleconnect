import { MissingEndpointError, request } from "./http";
import {
  deleteCollectionItem,
  getDashboardSnapshot,
  insertCollectionItem,
  readCollection,
  updateCollectionItem
} from "./mockDb";

const BACKOFFICE_API_BASE =
  import.meta.env.VITE_BACKOFFICE_API_BASE ?? "http://localhost:8080/api";

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
  const response = await request(`${BACKOFFICE_API_BASE}/admin/users`);
  return (response.items ?? []).map(normalizeUser);
}

async function readPrestationsFromApi(filters = {}) {
  const response = await request(`${BACKOFFICE_API_BASE}/annonces${buildQuery({ type: filters.type })}`);
  return (response.items ?? []).map(normalizePrestation);
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
        request(`${BACKOFFICE_API_BASE}/admin/metrics`),
        readUsersFromApi(),
        readPrestationsFromApi()
      ]);

      return buildDashboard({
        users,
        prestations,
        categories: await readCategoriesLocal(),
        events: await readEventsLocal(),
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
    return insertCollectionItem(
      "users",
      {
        ...payload,
        createdAt: payload.createdAt ?? new Date().toISOString().slice(0, 10)
      },
      "u"
    );
  },

  async updateUser(id, payload) {
    return updateCollectionItem("users", id, payload);
  },

  async toggleUserStatus(id) {
    const users = await readUsersLocal();
    const user = users.find((entry) => entry.id === id);
    if (!user) {
      throw new MissingEndpointError("utilisateur.toggle", "Utilisateur introuvable en local.");
    }
    const nextStatus = user.status === "active" ? "inactive" : "active";
    return updateCollectionItem("users", id, { status: nextStatus });
  },

  async deleteUser(id) {
    return deleteCollectionItem("users", id);
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
      price: Number(payload.price ?? 0)
    };

    try {
      return await request(`${BACKOFFICE_API_BASE}/annonces`, {
        method: "POST",
        body: JSON.stringify(requestPayload)
      });
    } catch {
      return insertCollectionItem(
        "prestations",
        {
          ...requestPayload,
          status: payload.status ?? "draft",
          provider: payload.provider ?? "Equipe locale",
          createdAt: new Date().toISOString().slice(0, 10)
        },
        "p"
      );
    }
  },

  async updatePrestation(id, payload) {
    return updateCollectionItem("prestations", id, payload);
  },

  async deletePrestation(id) {
    return deleteCollectionItem("prestations", id);
  },

  async listCategories(filters = {}) {
    const items = await readCategoriesLocal();
    return paginate(filterCategories(items, filters), filters.page, filters.pageSize);
  },

  async createCategory(payload) {
    return insertCollectionItem("categories", payload, "c");
  },

  async updateCategory(id, payload) {
    return updateCollectionItem("categories", id, payload);
  },

  async deleteCategory(id) {
    return deleteCollectionItem("categories", id);
  },

  async listEvents(filters = {}) {
    const items = await readEventsLocal();
    return paginate(filterEvents(items, filters), filters.page, filters.pageSize);
  },

  async createEvent(payload) {
    return insertCollectionItem("events", payload, "e");
  },

  async updateEvent(id, payload) {
    return updateCollectionItem("events", id, payload);
  },

  async deleteEvent(id) {
    return deleteCollectionItem("events", id);
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
