<template>
    <div class="layout-wrapper public-dashboard">
        <header class="content-header">
            <div class="header-left">
                <p class="sidebar__category2">ACCUEIL > ESPACE SALARIÉ</p>
                <h1 class="hero-title1">Bonjour {{ prenom }} 👋</h1>
                <p class="classic-text">
                    Voici un résumé de votre activité d'animateur sur
                    UpcycleConnect
                </p>
            </div>
            <div class="header-actions-group">
                <router-link
                    to="/profil/createFormation"
                    class="btn-main-action"
                >
                    + Créer une formation
                </router-link>
            </div>
        </header>

        <div class="stats-grid dashboard-stats">
            <div class="card card--score">
                <p class="tag-score">MON ACTIVITÉ</p>
                <div class="score-value">
                    {{ stats.total_inscrits || 0 }} <span>inscrits</span>
                </div>
                <p class="score-level">
                    {{ stats.formations_actives || 0 }} formation(s) active(s)
                </p>
                <div class="score-footer">
                    <div class="mini-stat">
                        <strong>{{ formations.length }}</strong
                        ><br />Formations
                    </div>
                    <div class="mini-stat">
                        <strong>{{ tips.length }}</strong
                        ><br />Conseils publiés
                    </div>
                    <div class="mini-stat">
                        <strong>{{ evenements.length }}</strong
                        ><br />Événements
                    </div>
                </div>
            </div>

            <div class="card card--white">
                <div class="card-num">{{ formationsEnAttenteCount }}</div>
                <p class="text-dm">Formations en attente</p>
                <span class="badge badge--orange">
                    {{
                        formationsEnAttenteCount > 0
                            ? "Validation requise"
                            : "Tout est validé"
                    }}
                </span>
            </div>

            <div class="card card--white">
                <div class="card-num2">{{ moderationStats.signales || 0 }}</div>
                <p class="text-dm">Messages signalés</p>
                <span
                    class="badge"
                    :class="
                        moderationStats.signales > 0
                            ? 'badge--red'
                            : 'badge--green'
                    "
                >
                    {{
                        moderationStats.signales > 0
                            ? "À traiter"
                            : "Aucun signalement"
                    }}
                </span>
            </div>
        </div>

        <div class="section-container">
            <div class="section-header">
                <h2>Mes dernières formations</h2>
                <div
                    class="header-actions"
                    style="display: flex; gap: 10px; align-items: center"
                >
                    <button
                        v-if="formations.length > 0"
                        class="btn-secondary"
                        @click="goToAllFormations"
                    >
                        Voir plus ({{ formations.length }})
                    </button>
                    <router-link
                        to="/profil/createFormation"
                        class="btn-main-action1"
                    >
                        + Nouvelle formation
                    </router-link>
                </div>
            </div>
            <table class="data-table">
                <thead>
                    <tr>
                        <th>FORMATION</th>
                        <th>TYPE</th>
                        <th>VALIDATION</th>
                        <th>STATUT</th>
                        <th>INSCRITS</th>
                        <th>DATE</th>
                        <th>ACTIONS</th>
                    </tr>
                </thead>
                <tbody>
                    <tr
                        v-for="formation in formations.slice(0, 3)"
                        :key="formation.id"
                    >
                        <td>
                            <strong>{{ formation.titre }}</strong
                            ><br />
                            <small class="table-subtext">{{
                                formation.ville
                            }}</small>
                        </td>
                        <td>
                            <span class="tag-formation">FORMATION</span>
                        </td>
                        <td>
                            <span
                                :class="
                                    formation.est_valide === 'Valide'
                                        ? 'status-valid'
                                        : 'status-pending'
                                "
                            >
                                {{
                                    formation.est_valide === "Valide"
                                        ? "APPROUVÉE"
                                        : formation.est_valide === "Refuse"
                                          ? "REFUSÉE"
                                          : "EN ATTENTE"
                                }}
                            </span>
                        </td>
                        <td>
                            <span class="status-logistique">{{
                                formation.statut
                            }}</span>
                        </td>
                        <td>
                            {{ formation.nb_inscrit || 0 }} /
                            {{ formation.capacite_max }}
                        </td>
                        <td>{{ formatDate(formation.date_debut) }}</td>
                        <td class="actions-cell">
                            <button
                                class="btn-view"
                                @click="goToFormation(formation.id)"
                            >
                                Voir
                            </button>
                            <button
                                v-if="formation.est_valide !== 'Valide'"
                                class="btn-modify"
                                @click="goToModifyFormation(formation.id)"
                            >
                                Modifier
                            </button>
                            <button
                                class="btn-remove"
                                @click="removeFormation(formation.id)"
                            >
                                Supprimer
                            </button>
                        </td>
                    </tr>
                </tbody>
            </table>
            <div v-if="formations.length === 0" class="state-card">
                Vous n'avez pas encore créé de formation.
            </div>
        </div>

        <div class="section-container">
            <div class="section-header">
                <h2>Mes événements organisés</h2>
                <div
                    class="header-actions"
                    style="display: flex; gap: 10px; align-items: center"
                >
                    <button
                        v-if="evenements.length > 0"
                        class="btn-secondary"
                        @click="goToAllEvenements"
                    >
                        Voir plus ({{ evenements.length }})
                    </button>
                    <router-link
                        to="/profil/createEvenement"
                        class="btn-main-action1"
                    >
                        + Nouvel événement
                    </router-link>
                </div>
            </div>
            <table class="data-table">
                <thead>
                    <tr>
                        <th>ÉVÉNEMENT</th>
                        <th>TYPE</th>
                        <th>LIEU</th>
                        <th>DATE</th>
                        <th>ACTIONS</th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-for="evt in evenements.slice(0, 3)" :key="evt.id">
                        <td>
                            <strong>{{ evt.titre }}</strong>
                        </td>
                        <td>
                            <span class="tag-event">ÉVÉNEMENT</span>
                        </td>
                        <td>{{ evt.ville }}</td>
                        <td>{{ formatDate(evt.date_evenement) }}</td>
                        <td class="actions-cell">
                            <button
                                class="btn-view"
                                @click="goToEvenement(evt.id)"
                            >
                                Voir
                            </button>
                            <button
                                class="btn-modify"
                                @click="goToModifyEvenement(evt.id)"
                            >
                                Modifier
                            </button>
                            <button
                                class="btn-remove"
                                @click="removeEvenement(evt.id)"
                            >
                                Supprimer
                            </button>
                        </td>
                    </tr>
                </tbody>
            </table>
            <div v-if="evenements.length === 0" class="state-card">
                Vous n'avez pas encore organisé d'événement.
            </div>
        </div>

        <div class="section-container planning-section">
            <div class="section-header">
                <div>
                    <h2>Mon planning ({{ planningWeekLabel }})</h2>
                    <p class="classic-text">
                        Vos formations et événements à venir.
                    </p>
                </div>
                <button class="btn-secondary" @click.stop="goToFullPlanning">
                    Vue mensuelle
                </button>
            </div>

            <div class="planning-week">
                <div
                    class="planning-week__weekday"
                    v-for="day in compactWeekDays"
                    :key="day.key"
                >
                    <span>{{ day.weekLabel }}</span>
                    <strong>{{ day.dayLabel }}</strong>
                </div>
                <div
                    v-for="day in compactWeekDays"
                    :key="`${day.key}-card`"
                    class="planning-week__day"
                    :class="{ 'is-today': day.isToday }"
                >
                    <div class="planning-week__head">{{ day.dayNumber }}</div>
                    <div class="planning-week__entries">
                        <div
                            v-for="entry in day.entries.slice(0, 2)"
                            :key="entry.id"
                            class="planning-week__entry"
                            :class="`planning-entry--${entry.kind}`"
                        >
                            <strong>{{ entry.title }}</strong>
                            <span>{{ entry.timeLabel }}</span>
                        </div>
                    </div>
                </div>
            </div>

            <div
                v-if="calendarEntries.length === 0"
                class="state-card"
                style="margin-top: 1rem"
            >
                Aucune formation ou événement prévu cette semaine.
            </div>
        </div>

        <div class="end-grid">
            <div class="dashboard-bottom-card" v-if="tips.length > 0">
                <div class="card-top-content">
                    <span class="tag-light-green">💡 Mon dernier conseil</span>
                    <h2 class="mod-title">{{ tips[0].titre }}</h2>
                    <p class="classic-text">
                        {{
                            tips[0].description ||
                            "Découvrez ce conseil publié sur la plateforme."
                        }}
                    </p>
                </div>

                <div class="card-bottom-actions">
                    <router-link
                        :to="{
                            name: 'conseil-detail',
                            params: { id: tips[0].id },
                        }"
                        class="btn-action-light-green"
                    >
                        Lire la suite →
                    </router-link>
                    <button class="btn-secondary" @click="goToAllTips">
                        Voir tous ({{ tips.length }})
                    </button>
                </div>
            </div>

            <div class="dashboard-bottom-card" v-else>
                <div class="card-top-content">
                    <span class="tag-light-green">💡 Conseils</span>
                    <h2 class="mod-title">Aucun conseil publié</h2>
                    <p class="classic-text">
                        Partagez votre expertise avec la communauté en créant
                        votre premier conseil !
                    </p>
                </div>
                <div class="card-bottom-actions">
                    <router-link
                        to="/profil/createTip"
                        class="btn-action-light-green"
                    >
                        + Publier un conseil
                    </router-link>
                </div>
            </div>

            <div class="dashboard-bottom-card moderation-custom-layout">
                <div class="mod-header">
                    <div class="mod-header-text">
                        <span class="tag-notif-yellow"
                            >🔔 Modération du forum</span
                        >
                        <h2 class="mod-title">Vue d'ensemble</h2>
                    </div>
                    <router-link
                        to="/profil/forum"
                        class="btn-action-light-green"
                    >
                        Gérer le forum →
                    </router-link>
                </div>

                <div class="mod-stats-grid">
                    <div
                        class="mod-stat-box"
                        :class="{ 'is-danger': moderationStats.signales > 0 }"
                    >
                        <div class="stat-icon">
                            <TriangleAlert :size="24" stroke-width="2.5" />
                        </div>
                        <div class="stat-content">
                            <div class="stat-val">
                                {{ moderationStats.signales || 0 }}
                            </div>
                            <div class="stat-lbl">Signalements</div>
                        </div>
                    </div>

                    <div class="mod-stat-box">
                        <div class="stat-icon">
                            <MessageCircle :size="24" stroke-width="2.5" />
                        </div>
                        <div class="stat-content">
                            <div class="stat-val">
                                {{ moderationStats.discussions || 0 }}
                            </div>
                            <div class="stat-lbl">Discussions actives</div>
                        </div>
                    </div>

                    <div class="mod-stat-box">
                        <div class="stat-icon">
                            <Ban :size="24" stroke-width="2.5" />
                        </div>
                        <div class="stat-content">
                            <div class="stat-val">
                                {{ moderationStats.bannis || 0 }}
                            </div>
                            <div class="stat-lbl">Comptes bannis</div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup>
import { computed, onMounted, ref } from "vue";
import { useRouter } from "vue-router";
import { TriangleAlert, MessageCircle, Ban } from "lucide-vue-next";

const router = useRouter();
const API_URL = "http://localhost:8081";

const prenom = ref(sessionStorage.getItem("userPrenom") || "Invité");
const formations = ref([]);
const evenements = ref([]);
const tips = ref([]);
const calendarEntries = ref([]);

const stats = ref({
    total_inscrits: 0,
    formations_actives: 0,
});

const moderationStats = ref({
    signales: 0,
    discussions: 0,
    bannis: 0,
});

const getLocalISODate = (date) => {
    const d = new Date(date);
    const year = d.getFullYear();
    const month = String(d.getMonth() + 1).padStart(2, "0");
    const day = String(d.getDate()).padStart(2, "0");
    return `${year}-${month}-${day}`;
};

const formationsEnAttenteCount = computed(
    () => formations.value.filter((f) => f.est_valide === "En attente").length,
);

const currentWeekStart = computed(() => {
    const now = new Date();
    const offset = (now.getDay() + 6) % 7;
    const monday = new Date(now);
    monday.setHours(0, 0, 0, 0);
    monday.setDate(now.getDate() - offset);
    return monday;
});

const planningWeekLabel = computed(() =>
    currentWeekStart.value.toLocaleDateString("fr-FR", {
        day: "numeric",
        month: "long",
    }),
);

const weekDays = ["Lun", "Mar", "Mer", "Jeu", "Ven", "Sam", "Dim"];

const compactWeekDays = computed(() =>
    Array.from({ length: 7 }, (_, index) => {
        const date = new Date(currentWeekStart.value);
        date.setDate(currentWeekStart.value.getDate() + index);
        const isoDate = getLocalISODate(date);

        return {
            key: isoDate,
            weekLabel: weekDays[index].toUpperCase(),
            dayLabel: date.toLocaleDateString("fr-FR", { day: "numeric" }),
            dayNumber: date.getDate(),
            isToday: isoDate === getLocalISODate(new Date()),
            entries: calendarEntries.value.filter((e) => e.date === isoDate),
        };
    }),
);

const formatDate = (val) => {
    if (!val) return "NULL";
    const date = new Date(val);
    return isNaN(date.getTime())
        ? "NULL"
        : new Intl.DateTimeFormat("fr-FR", {
              day: "2-digit",
              month: "short",
              year: "numeric",
          }).format(date);
};

const goToFormation = (id) =>
    router.push({ name: "formation-detail", params: { id } });
const goToModifyFormation = (id) =>
    router.push({ name: "modify-formation", params: { id } });
const goToEvenement = (id) =>
    router.push({ name: "evenement-detail", params: { id } });
const goToModifyEvenement = (id) =>
    router.push({ name: "modify-evenement", params: { id } });
const goToFullPlanning = () => router.push("/profil/planning");

const goToAllFormations = () => router.push("/profil/formations");
const goToAllEvenements = () => router.push("/profil/evenements");
const goToAllTips = () => router.push("/profil/tips");

const removeFormation = async (id) => {
    if (!confirm("Voulez-vous vraiment supprimer cette formation ?")) return;
    try {
        const token = sessionStorage.getItem("userToken") || "";
        const res = await fetch(`${API_URL}/formations/${id}`, {
            method: "DELETE",
            headers: { Authorization: `Bearer ${token}` },
        });
        if (res.ok) {
            formations.value = formations.value.filter((f) => f.id !== id);
        }
    } catch (e) {}
};

const removeEvenement = async (id) => {
    if (!confirm("Voulez-vous vraiment supprimer cet événement ?")) return;
    try {
        const token = sessionStorage.getItem("userToken") || "";
        const res = await fetch(`${API_URL}/evenements/${id}`, {
            method: "DELETE",
            headers: { Authorization: `Bearer ${token}` },
        });
        if (res.ok) {
            evenements.value = evenements.value.filter((e) => e.id !== id);
        }
    } catch (e) {}
};

onMounted(async () => {
    const rawId = sessionStorage.getItem("userId");
    if (!rawId || rawId === "null" || rawId === "undefined") return;

    const id = parseInt(rawId, 10);
    if (isNaN(id)) return;

    const headers = {
        Authorization: `Bearer ${sessionStorage.getItem("userToken") || ""}`,
    };

    try {
        const res = await fetch(`${API_URL}/formations`, { headers });
        if (res.ok) {
            const data = await res.json();
            const all = Array.isArray(data) ? data : data.formations || [];
            formations.value = all.filter((f) => f.id_formateur === id);
            stats.value.formations_actives = formations.value.filter(
                (f) => f.statut === "Ouvert",
            ).length;
            stats.value.total_inscrits = formations.value.reduce(
                (sum, f) => sum + (f.nb_inscrit || 0),
                0,
            );
        }
    } catch (e) {}

    try {
        const res = await fetch(`${API_URL}/evenements`, { headers });
        if (res.ok) {
            const data = await res.json();
            const all = Array.isArray(data) ? data : data.evenements || [];
            evenements.value = all.filter((e) => e.id_createur === id);
        }
    } catch (e) {}

    try {
        const res = await fetch(`${API_URL}/tips`, { headers });
        if (res.ok) {
            const data = await res.json();
            const all = Array.isArray(data) ? data : data.tips || [];
            tips.value = all.filter((t) => t.id_createur === id).reverse();
        }
    } catch (e) {}

    try {
        const [signalesRes, topicsRes, bannedRes] = await Promise.all([
            fetch(`${API_URL}/api/moderation/forums/signales`, { headers }),
            fetch(`${API_URL}/api/moderation/topics`, { headers }),
            fetch(`${API_URL}/api/moderation/users/banned`, { headers }),
        ]);

        const signales = signalesRes.ok ? await signalesRes.json() : [];
        const topics = topicsRes.ok ? await topicsRes.json() : [];
        const banned = bannedRes.ok ? await bannedRes.json() : [];

        moderationStats.value = {
            signales: Array.isArray(signales) ? signales.length : 0,
            discussions: Array.isArray(topics) ? topics.length : 0,
            bannis: Array.isArray(banned) ? banned.length : 0,
        };
    } catch (e) {}

    const mapEntry = (arr, kind, dateKey) =>
        (arr || [])
            .map((item) => {
                const parsed = new Date(item[dateKey]);
                if (isNaN(parsed.getTime())) return null;
                return {
                    id: `${kind}-${item.id}`,
                    kind,
                    title: item.titre,
                    date: getLocalISODate(parsed),
                    timeLabel: parsed.toLocaleTimeString("fr-FR", {
                        hour: "2-digit",
                        minute: "2-digit",
                    }),
                };
            })
            .filter(Boolean);

    calendarEntries.value = [
        ...mapEntry(formations.value, "formation", "date_debut"),
        ...mapEntry(evenements.value, "event", "date_evenement"),
    ].sort((a, b) => a.date.localeCompare(b.date));
});
</script>

<style scoped>
.public-dashboard {
    min-height: 100vh;
    padding: 20px;
    background: var(--bg-light, #f7f9f7);
}

.content-header {
    display: flex;
    justify-content: space-between;
    align-items: flex-end;
    margin-bottom: 2rem;
}

.header-actions-group {
    display: flex;
    gap: 0.8rem;
}

.section-container {
    background: white;
    padding: 2rem;
    border-radius: 16px;
    border: 1px solid #eee;
    margin-bottom: 2rem;
}

.section-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1rem;
}

.end-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
    gap: 1.5rem;
    margin-bottom: 2rem;
    align-items: stretch;
}

@media (min-width: 921px) {
    .end-grid {
        grid-template-columns: 1fr 1fr;
    }
}

.dashboard-bottom-card {
    background: white;
    padding: 2rem;
    border-radius: 16px;
    border: 1px solid #eee;
    display: flex;
    flex-direction: column;
    justify-content: space-between;
    height: 100%;
}

.card-top-content {
    display: flex;
    flex-direction: column;
    align-items: flex-start;
    gap: 12px;
}

.card-bottom-actions {
    display: flex;
    gap: 10px;
    margin-top: 2rem;
}

.btn-main-action,
.btn-main-action1 {
    display: inline-flex;
    align-items: center;
    background: #2d7a4f;
    color: white;
    padding: 10px 20px;
    border-radius: 10px;
    text-decoration: none;
    font-weight: bold;
    border: none;
    cursor: pointer;
}

.btn-secondary {
    padding: 8px 16px;
    border-radius: 8px;
    border: 1px solid #ddd;
    background: white;
    cursor: pointer;
    font-weight: 600;
}

.btn-action-light-green {
    background: #eaf4ed;
    color: #1b4d31;
    padding: 10px 20px;
    border-radius: 8px;
    font-weight: 700;
    font-size: 0.95rem;
    text-decoration: none;
    transition: all 0.2s ease;
    display: inline-block;
}

.btn-action-light-green:hover {
    background: #d1e7dd;
    transform: translateY(-2px);
}

.btn-view {
    background: transparent;
    border: 1px solid #ddd;
    padding: 6px 12px;
    border-radius: 6px;
    cursor: pointer;
    font-weight: bold;
}

.btn-view:hover {
    background: #f0f0f0;
}

.btn-modify {
    background: #fff3cd;
    border: 1px solid #ffeeba;
    color: #856404;
    padding: 6px 12px;
    border-radius: 6px;
    cursor: pointer;
    font-weight: bold;
    transition: all 0.2s;
}

.btn-modify:hover {
    background: #ffe8a1;
    color: #664d03;
}

.btn-remove {
    background: #ffe5e5;
    border: 1px solid #ffcccc;
    color: #d32f2f;
    padding: 6px 12px;
    border-radius: 6px;
    cursor: pointer;
    font-weight: bold;
    transition: all 0.2s;
}

.btn-remove:hover {
    background: #ffcccc;
    color: #b71c1c;
}

.tag-light-green {
    background: #eaf4ed;
    color: #2d7a4f;
    padding: 6px 16px;
    border-radius: 20px;
    font-size: 0.85rem;
    font-weight: 800;
}

.tag-notif-yellow {
    background: #fef3c7;
    color: #92400e;
    padding: 6px 16px;
    border-radius: 20px;
    font-size: 0.85rem;
    font-weight: 800;
}

.mod-title {
    font-size: 2rem;
    font-weight: 800;
    color: #1a1a1a;
    margin: 0;
}

.classic-text {
    color: #555;
    font-size: 0.95rem;
    line-height: 1.5;
    margin: 0;
}

.stats-grid {
    display: grid;
    gap: 1.5rem;
    margin-bottom: 2rem;
}

.dashboard-stats {
    grid-template-columns: 1.4fr 0.8fr 0.8fr;
}

@media (max-width: 920px) {
    .dashboard-stats {
        grid-template-columns: 1fr;
    }
}

.card {
    background: white;
    padding: 1.5rem;
    border-radius: 12px;
    border: 1px solid #eee;
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.02);
}

.card--score {
    background: #2d7a4f;
    color: white;
    border: none;
}

.score-value {
    font-size: 2.8rem;
    font-weight: 900;
    line-height: 1;
    margin-bottom: 0.5rem;
}

.score-value span {
    font-size: 1.2rem;
    font-weight: normal;
    opacity: 0.8;
}

.score-level {
    font-size: 1.1rem;
    margin-bottom: 1.5rem;
    opacity: 0.9;
}

.score-footer {
    display: flex;
    justify-content: space-between;
    border-top: 1px solid rgba(255, 255, 255, 0.2);
    padding-top: 1rem;
    margin-top: 1rem;
}

.mini-stat strong {
    font-size: 1.1rem;
    color: white;
    display: block;
}

.card-num,
.card-num2 {
    font-size: 2.5rem;
    font-weight: bold;
    color: #1a1a1a;
    line-height: 1;
    margin-bottom: 0.5rem;
}

.tag-score {
    background: rgba(255, 255, 255, 0.2);
    color: white;
    padding: 4px 10px;
    border-radius: 6px;
    font-size: 0.75rem;
    font-weight: 800;
    display: inline-block;
    margin-bottom: 1rem;
}

.tag-formation {
    background: #eaf4ed;
    color: #2d7a4f;
    padding: 4px 10px;
    border-radius: 6px;
    font-size: 0.75rem;
    font-weight: bold;
    text-transform: uppercase;
}

.tag-event {
    background: #fff3cd;
    color: #856404;
    padding: 4px 10px;
    border-radius: 6px;
    font-size: 0.75rem;
    font-weight: bold;
    text-transform: uppercase;
}

.status-valid {
    background: #eaf4ed;
    color: #2d7a4f;
    padding: 4px 10px;
    border-radius: 6px;
    font-size: 0.75rem;
    font-weight: bold;
}

.status-pending {
    background: #fff3cd;
    color: #856404;
    padding: 4px 10px;
    border-radius: 6px;
    font-size: 0.75rem;
    font-weight: bold;
}

.status-logistique {
    background: #f0f0f0;
    color: #555;
    padding: 4px 10px;
    border-radius: 6px;
    font-size: 0.75rem;
    font-weight: bold;
    text-transform: uppercase;
}

.badge {
    padding: 4px 10px;
    border-radius: 20px;
    font-size: 0.8rem;
    font-weight: bold;
    display: inline-block;
}

.badge--green {
    background: #eaf4ed;
    color: #2d7a4f;
}

.badge--orange {
    background: #fff3cd;
    color: #856404;
}

.badge--red {
    background: #fee2e2;
    color: #b91c1c;
}

.data-table {
    width: 100%;
    border-collapse: collapse;
}

.data-table th {
    text-align: left;
    font-size: 0.75rem;
    color: #999;
    text-transform: uppercase;
    padding: 10px;
    border-bottom: 2px solid #f0f0f0;
}

.data-table td {
    padding: 12px 10px;
    border-bottom: 1px solid #f5f5f5;
    font-size: 0.9rem;
}

.table-subtext {
    display: block;
    margin-top: 4px;
    color: var(--text-grey, #666);
    font-size: 0.78rem;
}

.actions-cell {
    display: flex;
    gap: 6px;
    flex-wrap: wrap;
}

.state-card {
    border: 1px dashed #cfe0d4;
    border-radius: 14px;
    padding: 26px;
    color: var(--text-grey, #666);
    background: #fbfdfb;
    text-align: center;
}

.planning-week {
    display: grid;
    grid-template-columns: repeat(7, 1fr);
    gap: 1rem;
}

.planning-week__weekday {
    text-align: center;
    font-weight: 700;
    color: #666;
    padding-bottom: 5px;
    display: flex;
    flex-direction: column;
}

.planning-week__day {
    background: #fbfdfb;
    border: 1px solid #e0e0e0;
    border-radius: 12px;
    min-height: 100px;
    padding: 10px;
    display: flex;
    flex-direction: column;
}

.planning-week__day.is-today {
    border: 2px solid #2d7a4f;
    background: #eaf4ed;
}

.planning-week__head {
    font-weight: bold;
    font-size: 1.1rem;
    margin-bottom: 8px;
    color: #1a1a1a;
}

.planning-week__entry {
    padding: 6px;
    border-radius: 6px;
    font-size: 0.75rem;
    margin-bottom: 4px;
    display: flex;
    flex-direction: column;
}

.planning-entry--formation {
    background: #d1e7dd;
    color: #0f5132;
}

.planning-entry--event {
    background: #fff3cd;
    color: #856404;
}

.moderation-custom-layout {
    display: flex;
    flex-direction: column;
    justify-content: space-between;
    height: 100%;
    gap: 2rem;
}

.mod-header {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    flex-wrap: wrap;
    gap: 1rem;
}

.mod-header-text {
    display: flex;
    flex-direction: column;
    align-items: flex-start;
    gap: 0.5rem;
}

.mod-stats-grid {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 15px;
}

@media (max-width: 650px) {
    .mod-stats-grid {
        grid-template-columns: 1fr;
    }
}

.mod-stat-box {
    background: #fafafa;
    border: 1px solid #f0f0f0;
    border-radius: 12px;
    padding: 1.2rem;
    display: flex;
    align-items: center;
    gap: 12px;
    transition: all 0.2s ease;
}

.stat-icon {
    color: #2d7a4f;
    background: white;
    width: 45px;
    height: 45px;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 10px;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.02);
}

.stat-content {
    display: flex;
    flex-direction: column;
}

.stat-val {
    font-size: 1.5rem;
    font-weight: 900;
    color: #1a1a1a;
    line-height: 1.1;
}

.stat-lbl {
    font-size: 0.75rem;
    font-weight: 700;
    color: #888;
    text-transform: uppercase;
    margin-top: 4px;
}

.mod-stat-box.is-danger {
    background: #fff5f5;
    border-color: #fecaca;
}

.mod-stat-box.is-danger .stat-val {
    color: #b91c1c;
}

.mod-stat-box.is-danger .stat-lbl {
    color: #991b1b;
}

.mod-stat-box.is-danger .stat-icon {
    background: #fee2e2;
    color: #b91c1c;
}
</style>
