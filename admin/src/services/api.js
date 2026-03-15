import { MissingEndpointError, request } from "./http";

const BACKOFFICE_API_BASE =
  import.meta.env.VITE_BACKOFFICE_API_BASE ?? "http://localhost:8080/api";
const GO_API_BASE = import.meta.env.VITE_GO_API_BASE ?? "http://localhost:8080/api-go";

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
  const fullName = row.name ?? `${firstName} ${lastName}`.trim();
  return {
    id: row.id,
    firstName,
    lastName,
    fullName,
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

export const capabilities = {
  dashboard: {
    metrics: true
  },
  users: {
    list: true,
    detail: false,
    create: false,
    update: false,
    toggle: false,
    delete: false
  },
  prestations: {
    list: true,
    create: true,
    update: false,
    delete: false
  },
  categories: {
    list: false,
    create: false,
    update: false,
    delete: false
  },
  events: {
    list: false,
    create: false,
    update: false,
    delete: false
  }
};

export const adminApi = {
  async getDashboard() {
    const [metricsResponse, usersResponse, prestationsResponse] = await Promise.all([
      request(`${BACKOFFICE_API_BASE}/admin/metrics`),
      request(`${BACKOFFICE_API_BASE}/admin/users`),
      request(`${BACKOFFICE_API_BASE}/annonces`)
    ]);

    const users = (usersResponse.items ?? []).map(normalizeUser);
    const prestations = (prestationsResponse.items ?? []).map(normalizePrestation);
    const metrics = metricsResponse.metrics ?? {};

    const usersByRole = users.reduce((acc, user) => {
      const key = user.role.toLowerCase();
      acc[key] = (acc[key] ?? 0) + 1;
      return acc;
    }, {});

    return {
      stats: [
        {
          label: "Utilisateurs",
          value: metrics.users ?? users.length,
          tone: "green"
        },
        {
          label: "Prestations",
          value: metrics.annonces ?? prestations.length,
          tone: "teal"
        },
        {
          label: "Categories",
          value: null,
          tone: "sand",
          missing: true
        },
        {
          label: "Evenements",
          value: null,
          tone: "amber",
          missing: true
        }
      ],
      charts: {
        usersByRole: [
          { label: "Particuliers", value: usersByRole.particulier ?? 0 },
          { label: "Prestataires", value: usersByRole.prestataire ?? usersByRole.pro ?? 0 },
          { label: "Admins", value: usersByRole.admin ?? 0 }
        ],
        resources: [
          { label: "Prestations", value: prestations.length },
          { label: "Categories", value: 0, missing: true },
          { label: "Evenements", value: 0, missing: true }
        ]
      },
      quickNotes: [
        users.length === 0
          ? { tone: "amber", text: "Aucun utilisateur remonte depuis l'API admin/users." }
          : { tone: "green", text: `${users.length} utilisateurs remontent dans le backoffice.` },
        {
          tone: "amber",
          text: "Les categories et evenements n'ont pas encore d'endpoints admin exposes."
        },
        prestations.length > 0
          ? { tone: "green", text: `${prestations.length} prestations recuperables via /annonces.` }
          : { tone: "coral", text: "Aucune prestation n'a ete retournee par /annonces." }
      ],
      recentActivity: prestations.slice(0, 5).map((item) => ({
        id: item.id,
        title: item.title,
        subtitle: item.description || "Prestation enregistree"
      }))
    };
  },

  async listUsers(filters = {}) {
    const response = await request(`${BACKOFFICE_API_BASE}/admin/users`);
    const normalized = (response.items ?? []).map(normalizeUser);
    let filtered = normalized;

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

    return paginate(filtered, filters.page, filters.pageSize);
  },

  async getUser() {
    throw new MissingEndpointError(
      "utilisateur.detail",
      "Le backend actuel n'expose pas encore GET /api/admin/users/:id ni equivalent."
    );
  },

  async createUser() {
    throw new MissingEndpointError(
      "utilisateur.create",
      "Le backend actuel n'expose pas encore POST /api/admin/users ni equivalent exploitable via la gateway."
    );
  },

  async updateUser() {
    throw new MissingEndpointError(
      "utilisateur.update",
      "Le backend actuel n'expose pas encore PUT /api/admin/users/:id ni equivalent exploitable via la gateway."
    );
  },

  async toggleUserStatus() {
    throw new MissingEndpointError(
      "utilisateur.toggle",
      "Le backend actuel n'expose pas encore d'activation ou de desactivation d'utilisateur."
    );
  },

  async deleteUser() {
    throw new MissingEndpointError(
      "utilisateur.delete",
      "Le backend actuel n'expose pas encore DELETE /api/admin/users/:id ni equivalent exploitable via la gateway."
    );
  },

  async listPrestations(filters = {}) {
    const response = await request(`${BACKOFFICE_API_BASE}/annonces${buildQuery({ type: filters.type })}`);
    let items = (response.items ?? []).map(normalizePrestation);

    if (filters.search) {
      const needle = filters.search.toLowerCase();
      items = items.filter((item) =>
        [item.title, item.description, item.provider].join(" ").toLowerCase().includes(needle)
      );
    }

    if (filters.status) {
      items = items.filter((item) => item.status === filters.status);
    }

    if (filters.sortBy === "price-desc") {
      items = items.slice().sort((a, b) => b.price - a.price);
    } else if (filters.sortBy === "price-asc") {
      items = items.slice().sort((a, b) => a.price - b.price);
    } else {
      items = items.slice().sort((a, b) => a.title.localeCompare(b.title));
    }

    return paginate(items, filters.page, filters.pageSize);
  },

  async createPrestation(payload) {
    const requestPayload = {
      title: payload.title,
      description: payload.description,
      type: payload.type,
      price: Number(payload.price ?? 0)
    };

    return request(`${BACKOFFICE_API_BASE}/annonces`, {
      method: "POST",
      body: JSON.stringify(requestPayload)
    });
  },

  async updatePrestation() {
    throw new MissingEndpointError(
      "prestation.update",
      "Le backend actuel n'expose pas encore PUT /api/annonces/:id."
    );
  },

  async deletePrestation() {
    throw new MissingEndpointError(
      "prestation.delete",
      "Le backend actuel n'expose pas encore DELETE /api/annonces/:id."
    );
  },

  async listCategories() {
    throw new MissingEndpointError(
      "categorie.list",
      "Le backend actuel n'expose pas encore GET /api/categories."
    );
  },

  async createCategory() {
    throw new MissingEndpointError(
      "categorie.create",
      "Le backend actuel n'expose pas encore POST /api/categories."
    );
  },

  async updateCategory() {
    throw new MissingEndpointError(
      "categorie.update",
      "Le backend actuel n'expose pas encore PUT /api/categories/:id."
    );
  },

  async deleteCategory() {
    throw new MissingEndpointError(
      "categorie.delete",
      "Le backend actuel n'expose pas encore DELETE /api/categories/:id."
    );
  },

  async listEvents() {
    throw new MissingEndpointError(
      "evenement.list",
      "Le backend actuel n'expose pas encore GET /api/events."
    );
  },

  async createEvent() {
    throw new MissingEndpointError(
      "evenement.create",
      "Le backend actuel n'expose pas encore POST /api/events."
    );
  },

  async updateEvent() {
    throw new MissingEndpointError(
      "evenement.update",
      "Le backend actuel n'expose pas encore PUT /api/events/:id."
    );
  },

  async deleteEvent() {
    throw new MissingEndpointError(
      "evenement.delete",
      "Le backend actuel n'expose pas encore DELETE /api/events/:id."
    );
  },

  getCapabilities() {
    return capabilities;
  }
};
