const API = {
  backoffice: "/api",
  go: "/api-go"
};

const state = {
  lang: "fr",
  userId: "demo-user-1"
};

const translations = {
  fr: {
    subtitle: "Plateforme upcycling pour particuliers, pros et salaries."
  },
  en: {
    subtitle: "Upcycling platform for citizens, professionals and employees."
  }
};

function byId(id) {
  return document.getElementById(id);
}

function setFeedback(id, message) {
  byId(id).textContent = message;
}

async function getJSON(url, opts = {}) {
  const res = await fetch(url, {
    headers: { "Content-Type": "application/json" },
    ...opts
  });
  if (!res.ok) throw new Error(`HTTP ${res.status}`);
  return res.json();
}

function bindTabs() {
  document.querySelectorAll(".tab").forEach((btn) => {
    btn.addEventListener("click", () => {
      document.querySelectorAll(".tab").forEach((x) => x.classList.remove("active"));
      document.querySelectorAll(".panel").forEach((x) => x.classList.remove("active"));
      btn.classList.add("active");
      byId(btn.dataset.tab).classList.add("active");
    });
  });
}

function bindLangSwitch() {
  document.querySelectorAll(".lang-switch button").forEach((btn) => {
    btn.addEventListener("click", () => {
      state.lang = btn.dataset.lang;
      document.querySelectorAll(".lang-switch button").forEach((x) => x.classList.remove("active"));
      btn.classList.add("active");
      byId("subtitle").textContent = translations[state.lang].subtitle;
    });
  });
}

async function loadConseils() {
  const data = await getJSON(`${API.backoffice}/conseils`);
  const ul = byId("conseils-list");
  ul.innerHTML = "";
  data.items.forEach((item) => {
    const li = document.createElement("li");
    li.textContent = item.title;
    ul.appendChild(li);
  });
}

async function loadFormations() {
  const data = await getJSON(`${API.backoffice}/formations`);
  const wrap = byId("formations-list");
  wrap.innerHTML = "";
  data.items.forEach((item) => {
    const card = document.createElement("article");
    card.innerHTML = `
      <h4>${item.title}</h4>
      <p>${item.description}</p>
      <p><strong>${item.priceEUR} EUR</strong></p>
      <button data-id="${item.id}">Acheter (Stripe a connecter)</button>
    `;
    card.querySelector("button").addEventListener("click", () => {
      alert(`Paiement Stripe a brancher pour la formation ${item.id}`);
    });
    wrap.appendChild(card);
  });
}

async function loadMarket() {
  const data = await getJSON(`${API.backoffice}/annonces?type=vente`);
  const list = byId("market-list");
  list.innerHTML = "";
  data.items.forEach((item) => {
    const li = document.createElement("li");
    li.textContent = `${item.title} - ${item.price ?? 0} EUR`;
    list.appendChild(li);
  });
}

async function loadPlanning() {
  const data = await getJSON(`${API.backoffice}/planning?userId=${encodeURIComponent(state.userId)}`);
  const list = byId("planning-list");
  list.innerHTML = "";
  data.items.forEach((item) => {
    const li = document.createElement("li");
    li.textContent = `${item.date}: ${item.task}`;
    list.appendChild(li);
  });
}

function bindForms() {
  byId("annonce-form").addEventListener("submit", async (e) => {
    e.preventDefault();
    const fd = new FormData(e.target);
    const payload = Object.fromEntries(fd.entries());
    if (payload.type === "don") payload.price = 0;
    try {
      await getJSON(`${API.backoffice}/annonces`, {
        method: "POST",
        body: JSON.stringify(payload)
      });
      setFeedback("annonce-feedback", "Annonce publiee.");
      e.target.reset();
      loadMarket();
    } catch {
      setFeedback("annonce-feedback", "Erreur publication annonce.");
    }
  });

  byId("barcode-form").addEventListener("submit", async (e) => {
    e.preventDefault();
    const payload = Object.fromEntries(new FormData(e.target).entries());
    try {
      const validation = await getJSON(`${API.go}/barcode/validate`, {
        method: "POST",
        body: JSON.stringify(payload)
      });
      await getJSON(`${API.backoffice}/container-deposits`, {
        method: "POST",
        body: JSON.stringify({ ...payload, status: validation.status })
      });
      setFeedback("barcode-feedback", `Depot ${validation.status}.`);
      e.target.reset();
    } catch {
      setFeedback("barcode-feedback", "Erreur validation code barre.");
    }
  });

  byId("score-form").addEventListener("submit", async (e) => {
    e.preventDefault();
    const payload = Object.fromEntries(new FormData(e.target).entries());
    try {
      const data = await getJSON(`${API.go}/upcycling-score`, {
        method: "POST",
        body: JSON.stringify(payload)
      });
      setFeedback("score-feedback", `Score: ${data.score}/100 - CO2 evite: ${data.co2Kg} kg`);
    } catch {
      setFeedback("score-feedback", "Erreur calcul score.");
    }
  });

  byId("planning-form").addEventListener("submit", async (e) => {
    e.preventDefault();
    const payload = Object.fromEntries(new FormData(e.target).entries());
    payload.userId = state.userId;
    try {
      await getJSON(`${API.backoffice}/planning`, {
        method: "POST",
        body: JSON.stringify(payload)
      });
      e.target.reset();
      loadPlanning();
    } catch {
      setFeedback("score-feedback", "Erreur ajout planning.");
    }
  });

  byId("subscription-form").addEventListener("submit", async (e) => {
    e.preventDefault();
    const payload = Object.fromEntries(new FormData(e.target).entries());
    try {
      const data = await getJSON(`${API.backoffice}/pro/subscriptions`, {
        method: "POST",
        body: JSON.stringify(payload)
      });
      setFeedback("subscription-feedback", `Abonnement ${data.plan} actif.`);
    } catch {
      setFeedback("subscription-feedback", "Erreur abonnement.");
    }
  });

  byId("project-form").addEventListener("submit", async (e) => {
    e.preventDefault();
    const payload = Object.fromEntries(new FormData(e.target).entries());
    try {
      await getJSON(`${API.backoffice}/pro/projects`, {
        method: "POST",
        body: JSON.stringify(payload)
      });
      setFeedback("project-feedback", "Projet publie.");
      e.target.reset();
    } catch {
      setFeedback("project-feedback", "Erreur publication projet.");
    }
  });

  byId("employee-training-form").addEventListener("submit", async (e) => {
    e.preventDefault();
    const payload = Object.fromEntries(new FormData(e.target).entries());
    try {
      await getJSON(`${API.backoffice}/employee/trainings`, {
        method: "POST",
        body: JSON.stringify(payload)
      });
      setFeedback("employee-training-feedback", "Formation employee creee.");
      e.target.reset();
    } catch {
      setFeedback("employee-training-feedback", "Erreur creation formation.");
    }
  });

  byId("news-form").addEventListener("submit", async (e) => {
    e.preventDefault();
    const payload = Object.fromEntries(new FormData(e.target).entries());
    try {
      await getJSON(`${API.backoffice}/news`, {
        method: "POST",
        body: JSON.stringify(payload)
      });
      setFeedback("news-feedback", "News publiee.");
      e.target.reset();
    } catch {
      setFeedback("news-feedback", "Erreur publication news.");
    }
  });

  byId("moderation-form").addEventListener("submit", async (e) => {
    e.preventDefault();
    const payload = Object.fromEntries(new FormData(e.target).entries());
    try {
      await getJSON(`${API.backoffice}/employee/moderation`, {
        method: "POST",
        body: JSON.stringify(payload)
      });
      setFeedback("moderation-feedback", "Action de moderation enregistree.");
      e.target.reset();
    } catch {
      setFeedback("moderation-feedback", "Erreur moderation.");
    }
  });

  byId("load-admin-metrics").addEventListener("click", async () => {
    try {
      const data = await getJSON(`${API.backoffice}/admin/metrics`);
      byId("admin-output").textContent = JSON.stringify(data, null, 2);
    } catch {
      byId("admin-output").textContent = "Impossible de charger les metriques.";
    }
  });
}

function showTutorialOnce() {
  const key = "uc_tutorial_done";
  const modal = byId("tutorial-modal");
  if (!localStorage.getItem(key)) {
    modal.hidden = false;
  }
  byId("close-tutorial").addEventListener("click", () => {
    localStorage.setItem(key, "1");
    modal.hidden = true;
  });
}

async function init() {
  bindTabs();
  bindLangSwitch();
  bindForms();
  showTutorialOnce();

  await Promise.allSettled([
    loadConseils(),
    loadFormations(),
    loadMarket(),
    loadPlanning()
  ]);
}

init();
