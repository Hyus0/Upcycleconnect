<template>
    <div class="layout-wrapper public-dashboard">
        <header class="content-header">
            <div class="header-left">
                <p class="sidebar__category2">
                    {{ t.Home || "Accueil" }} >
                    {{
                        t.Dashboard
                            ? t.Dashboard.toUpperCase()
                            : "TABLEAU DE BORD"
                    }}
                </p>
                <h1 class="hero-title1">
                    {{ t.Hello || "Bonjour" }} {{ prenom }} 👋
                </h1>
                <p class="classic-text">
                    {{
                        t.SummaryActivity ||
                        "Voici un résumé de votre activité sur UpcycleConnect."
                    }}
                </p>
            </div>
            <router-link to="/profil/createAnnonce" class="btn-main-action">{{
                t.DeposeAnnonce || "Déposer une annonce"
            }}</router-link>
        </header>

        <div class="stats-grid">
            <div class="card card--score">
                <p class="tag-score">UPCYCLING SCORE</p>
                <div class="score-value">
                    {{ stats.total_points || 0 }} <span>pts</span>
                </div>
                <p class="score-level">
                    {{ t.Niveau || "Niveau" }} : {{ stats.niveau || "Novice" }}
                </p>
                <div class="score-footer">
                    <div class="mini-stat">
                        <strong>{{ stats.co2_total_evite_kg || 0 }} kg</strong
                        ><br />{{ t.Co2 || "CO2" }}
                    </div>
                    <div class="mini-stat">
                        <strong>{{ stats.nb_objets_recycles || 0 }}</strong
                        ><br />{{ t.Objets || "Objets" }}
                    </div>
                    <div class="mini-stat">
                        <strong
                            >{{ stats.ressources_economisees || 0 }} €</strong
                        ><br />{{ t.Economise || "Économisés" }}
                    </div>
                </div>
            </div>

            <div class="card">
                <div class="card-num">{{ projets.length }}</div>
                <p class="text-dm">Mes Projets</p>
                <span class="badge badge--green">En cours d'upcycling</span>
            </div>

            <div class="card">
                <div class="card-num" style="color: #d32f2f">
                    {{ achats.length }}
                </div>
                <p class="text-dm">Récupérations en attente</p>
                <span class="badge badge--orange">{{
                    achats.length > 0 ? "À récupérer" : "Aucun retrait"
                }}</span>
            </div>
        </div>

        <div class="section-container">
            <div class="section-header">
                <h2>Récupérations en attente</h2>
                <router-link
                    to="/profil/recuperations"
                    class="btn-secondary"
                    style="text-decoration: none"
                    >Voir tout</router-link
                >
            </div>
            <table class="data-table">
                <thead>
                    <tr>
                        <th>Objet</th>
                        <th>Conteneur / Site</th>
                        <th>Code de retrait</th>
                        <th>Statut</th>
                        <th class="text-right">Actions</th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-for="achat in achats.slice(0, 4)" :key="achat.id">
                        <td>
                            <strong>{{ achat.titre_annonce || "Objet" }}</strong
                            ><br />
                            <small class="table-subtext">{{
                                formatDate(achat.date_transaction)
                            }}</small>
                        </td>
                        <td>
                            {{ achat.nom_site || "En attente d'assignation" }}
                        </td>
                        <td>
                            <strong style="letter-spacing: 2px">{{
                                achat.code_retrait || "---"
                            }}</strong>
                        </td>
                        <td>
                            <span
                                class="badge"
                                :class="
                                    achat.statut === 'Récupéré'
                                        ? 'badge--green'
                                        : 'status-logistique'
                                "
                            >
                                {{ achat.statut || "À récupérer" }}
                            </span>
                        </td>
                        <td class="text-right">
                            <button
                                v-if="achat.statut !== 'Récupéré'"
                                class="btn-view"
                                @click="$router.push('/profil/depots')"
                            >
                                Suivre
                            </button>
                            <span v-else class="text-sm text-gray-500"
                                >Terminé</span
                            >
                        </td>
                    </tr>
                    <tr v-if="achats.length === 0">
                        <td colspan="5" class="text-center text-grey py-4">
                            Vous n'avez aucune récupération en attente.
                        </td>
                    </tr>
                </tbody>
            </table>
        </div>

        <div class="section-container">
            <div class="section-header">
                <h2>Mes Projets</h2>
                <router-link
                    to="/projets"
                    class="btn-secondary"
                    style="text-decoration: none"
                    >Voir le catalogue</router-link
                >
            </div>
            <table class="data-table">
                <thead>
                    <tr>
                        <th style="width: 40%">Projet</th>
                        <th>IMPACT CO2</th>
                        <th>ENGAGEMENT</th>
                        <th>STATUT</th>
                        <th>DATE</th>
                        <th>ACTIONS</th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-for="projet in projets.slice(0, 4)" :key="projet.id">
                        <td>
                            <strong>{{ projet.titre }}</strong
                            ><br />
                            <small class="text-truncate">{{
                                projet.description_courte ||
                                "Aucune description"
                            }}</small>
                        </td>
                        <td>
                            <span class="material-tag">
                                {{ projet.co2_evite_kg || 0 }} kg
                            </span>
                        </td>
                        <td>
                            <span class="stats-text"
                                >{{ projet.nb_vues || 0 }} vues,
                                {{ projet.nb_likes || 0 }} likes</span
                            ><br />
                        </td>
                        <td>
                            <span
                                :class="
                                    projet.visible_public
                                        ? 'status-valid'
                                        : 'status-neutral'
                                "
                            >
                                {{ projet.visible_public ? "PUBLIC" : "PRIVÉ" }}
                            </span>
                        </td>
                        <td>{{ formatDate(projet.date_creation) }}</td>
                        <td class="actions-cell">
                            <button
                                class="btn-view"
                                type="button"
                                @click="$router.push(`/projets/${projet.id}`)"
                            >
                                {{ t.Voir || "Voir" }}
                            </button>
                            <button
                                class="btn-modify"
                                type="button"
                                @click="
                                    $router.push(
                                        `/profil/modifyProjet/${projet.id}`,
                                    )
                                "
                            >
                                {{ t.Modifier || "Modifier" }}
                            </button>
                            <button
                                class="btn-remove"
                                type="button"
                                @click="removeProjet(projet.id)"
                            >
                                {{ t.Supprimer || "Retirer" }}
                            </button>
                        </td>
                    </tr>
                    <tr v-if="projets.length === 0">
                        <td colspan="6" class="text-center text-grey py-4">
                            Vous n'avez créé aucun projet.
                        </td>
                    </tr>
                </tbody>
            </table>
        </div>

        <div class="section-container">
            <div class="section-header">
                <div>
                    <h2>
                        {{ t.MySchedule || "Mon Planning" }} ({{
                            planningWeekLabel
                        }})
                    </h2>
                    <p class="classic-text">
                        {{
                            t.Schedule ||
                            "Vos prochains événements et formations."
                        }}
                    </p>
                </div>
                <button
                    class="btn-secondary"
                    @click="$router.push('/profil/planning')"
                >
                    {{ t.VueMensuelle || "Vue Mensuelle" }}
                </button>
            </div>

            <div class="planning-week">
                <div
                    class="planning-day"
                    v-for="day in compactWeekDays"
                    :key="day.key"
                    :class="{ 'is-today': day.isToday }"
                >
                    <div class="day-label">
                        {{ day.weekLabel }} <strong>{{ day.dayLabel }}</strong>
                    </div>
                    <div class="entries">
                        <div
                            v-for="entry in day.entries.slice(0, 2)"
                            :key="entry.id"
                            class="entry"
                            :class="`entry--${entry.kind}`"
                        >
                            <strong>{{ entry.title }}</strong>
                            <span>{{ entry.timeLabel }}</span>
                        </div>
                    </div>
                </div>
            </div>
            <div v-if="calendarEntries.length === 0" class="state-card mt-4">
                {{ t.PlanningInfo || "Rien de prévu." }}
            </div>
        </div>

        <div class="end-grid">
            <div class="section-container">
                <p class="badge badge--green mb-2">
                    💡 {{ t.Tips || "Conseil du jour" }}
                </p>
                <h2 v-if="tipDuJour">{{ tipDuJour.titre }}</h2>
                <h2 v-else>{{ t.WaitingTips || "En attente de conseil" }}</h2>
                <p class="mt-2">
                    {{
                        tipDuJour
                            ? tipDuJour.description
                            : t.TipsText || "Revenez plus tard."
                    }}
                </p>
                <router-link
                    v-if="tipDuJour"
                    :to="`/conseil/${tipDuJour.id}`"
                    class="btn-text-green mt-4"
                    style="text-decoration: none"
                >
                    {{ t.ReadNext || "Lire la suite" }} →
                </router-link>
            </div>

            <div class="section-container">
                <p class="badge badge--orange mb-2">
                    🔔 {{ t.Notification || "Notification" }}
                </p>
                <h2 v-if="latestNotification">
                    {{ latestNotification.titre }}
                </h2>
                <h2 v-else>
                    {{ t.NothingNotification || "Aucune notification" }}
                </h2>
                <p class="mt-2">
                    {{
                        latestNotification
                            ? latestNotification.message
                            : t.NotificationText || "Rien de nouveau."
                    }}
                </p>
                <button
                    v-if="latestNotification"
                    class="btn-text-green mt-4"
                    @click="handleNotificationClick(latestNotification)"
                >
                    {{ t.SeeAction || "Voir l'action" }} →
                </button>
            </div>
        </div>
    </div>
</template>

<script setup>
import { computed, onMounted, ref } from "vue";
import { useRouter } from "vue-router";
import { useTraduction } from "../../composables/useTraduction";

const router = useRouter();
const API_URL = "http://localhost:8081";
const { t } = useTraduction();

const prenom = ref(sessionStorage.getItem("userPrenom") || "Invité");
const achats = ref([]);
const projets = ref([]);
const calendarEntries = ref([]);
const tipDuJour = ref(null);
const latestNotification = ref(null);

const stats = ref({
    total_points: 0,
    niveau: "...",
    co2_total_evite_kg: 0,
    nb_objets_recycles: 0,
    ressources_economisees: 0,
});

const getLocalISODate = (date) => {
    const d = new Date(date);
    return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, "0")}-${String(d.getDate()).padStart(2, "0")}`;
};

const currentWeekStart = computed(() => {
    const now = new Date();
    now.setHours(0, 0, 0, 0);
    now.setDate(now.getDate() - ((now.getDay() + 6) % 7));
    return now;
});

const planningWeekLabel = computed(() =>
    currentWeekStart.value.toLocaleDateString("fr-FR", {
        day: "numeric",
        month: "long",
    }),
);

const weekDays = computed(() => [
    t.value.Lundi || "LUN",
    t.value.Mardi || "MAR",
    t.value.Mercredi || "MER",
    t.value.Jeudi || "JEU",
    t.value.Vendredi || "VEN",
    t.value.Samedi || "SAM",
    t.value.Dimanche || "DIM",
]);

const compactWeekDays = computed(() =>
    Array.from({ length: 7 }, (_, i) => {
        const date = new Date(currentWeekStart.value);
        date.setDate(date.getDate() + i);
        const iso = getLocalISODate(date);
        return {
            key: iso,
            weekLabel: weekDays.value[i].toUpperCase(),
            dayLabel: date.toLocaleDateString("fr-FR", { day: "numeric" }),
            isToday: iso === getLocalISODate(new Date()),
            entries: calendarEntries.value.filter((e) => e.date === iso),
        };
    }),
);

const formatDate = (val) => {
    if (!val) return "---";
    const d = new Date(val);
    return isNaN(d)
        ? "---"
        : new Intl.DateTimeFormat("fr-FR", {
              day: "2-digit",
              month: "short",
              year: "numeric",
          }).format(d);
};

const handleNotificationClick = async (notif) => {
    if (!notif) return;
    const token = sessionStorage.getItem("userToken") || "";
    const userId = sessionStorage.getItem("userId");
    try {
        await fetch(`${API_URL}/notifications/${notif.id}/read`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                Authorization: `Bearer ${token}`,
            },
            body: JSON.stringify({ user_id: parseInt(userId, 10) }),
        });
    } catch (e) {
        console.error(e);
    }
    router.push(
        notif.titre.toLowerCase().includes("disponible") ||
            notif.type === "Alerte"
            ? "/profil/depots"
            : "",
    );
};

const removeProjet = async (id) => {
    if (
        !confirm(
            t.value.ConfirmDeleteProjet ||
                "Voulez-vous vraiment supprimer ce projet ?",
        )
    )
        return;
    try {
        const token = sessionStorage.getItem("userToken") || "";
        const res = await fetch(`${API_URL}/projets/${id}`, {
            method: "DELETE",
            headers: { Authorization: `Bearer ${token}` },
        });
        if (!res.ok) throw new Error("Erreur suppression projet");
        projets.value = projets.value.filter((p) => p.id !== id);
    } catch (e) {
        console.error("Erreur suppression:", e);
        alert(
            t.value.ErrorDeleteProjet ||
                "Une erreur est survenue lors de la suppression.",
        );
    }
};

onMounted(() => {
    const id = parseInt(sessionStorage.getItem("userId"), 10);
    if (isNaN(id)) return;

    const headers = {
        Authorization: `Bearer ${sessionStorage.getItem("userToken") || ""}`,
    };
    const role = sessionStorage.getItem("userRole") || "Particulier";

    Promise.all([
        fetch(`${API_URL}/users/${id}/stats`, { headers })
            .then((r) => r.json())
            .then((d) => {
                stats.value = d;
                if (d.total_points)
                    sessionStorage.setItem("userScore", d.total_points);
            })
            .catch(() => {}),

        fetch(`${API_URL}/users/${id}/achats`, { headers })
            .then((r) => r.json())
            .then((d) => (achats.value = Array.isArray(d) ? d : []))
            .catch(() => {}),
        fetch(`${API_URL}/users/${id}/projets`, { headers })
            .then((r) => r.json())
            .then((d) => (projets.value = Array.isArray(d) ? d : []))
            .catch(() => {}),
        fetch(`${API_URL}/tips/role/${role}`)
            .then((r) => r.json())
            .then((d) => (tipDuJour.value = d))
            .catch(() => {}),
        fetch(`${API_URL}/users/${id}/notifications`, { headers })
            .then((r) => r.json())
            .then((d) => (latestNotification.value = d[0] || null))
            .catch(() => {}),
        fetch(`${API_URL}/user/planning/${id}`, { headers })
            .then((r) => r.json())
            .then((d) => {
                if (!d || Array.isArray(d)) return;
                const mapE = (arr, kind, key) =>
                    (arr || [])
                        .map((i) => {
                            const p = new Date(i[key]);
                            return isNaN(p)
                                ? null
                                : {
                                      id: `${kind}-${i.id}`,
                                      kind,
                                      title: i.titre,
                                      date: getLocalISODate(p),
                                      timeLabel: p.toLocaleTimeString("fr-FR", {
                                          hour: "2-digit",
                                          minute: "2-digit",
                                      }),
                                  };
                        })
                        .filter(Boolean);
                calendarEntries.value = [
                    ...mapE(d.formations, "formation", "date_debut"),
                    ...mapE(d.evenements, "event", "date_evenement"),
                ].sort((a, b) => a.date.localeCompare(b.date));
            })
            .catch(() => {}),
    ]);
});
</script>

<style scoped>
.public-dashboard {
    padding: 20px;
    background: #f7f9f7;
    min-height: 100vh;
}
.content-header {
    display: flex;
    justify-content: space-between;
    align-items: flex-end;
    margin-bottom: 2rem;
}
.hero-title1 {
    font-size: 2.2rem;
    font-weight: 900;
    color: #122018;
    margin: 0;
}
.sidebar__category2 {
    margin: 0;
    color: #a0ada7;
    font-family: "Space Mono", monospace;
    font-size: 0.65rem;
    text-transform: uppercase;
}
.classic-text {
    font-size: 0.95rem;
    color: #63746a;
    margin-top: 8px;
    line-height: 1.5;
}

.btn-main-action,
.btn-secondary,
.btn-text-green,
.btn-view {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    text-decoration: none;
    font-weight: bold;
    border-radius: 10px;
    transition: 0.2s;
}
.btn-main-action {
    background: #2d7a4f;
    color: white;
    padding: 10px 20px;
    border: none;
}
.btn-main-action:hover {
    background: #1b4d31;
}
.btn-secondary {
    padding: 8px 16px;
    border: 1px solid #ddd;
    background: white;
    color: #333;
}
.btn-secondary:hover {
    background: #f0f0f0;
}
.btn-view {
    background: transparent;
    border: 1px solid #ddd;
    padding: 6px 12px;
    border-radius: 6px;
    color: #333;
}
.btn-view:hover {
    background: #f0f0f0;
}
.btn-text-green {
    background: #eaf4ed;
    color: #2d7a4f;
    padding: 12px 20px;
    width: 100%;
    border: none;
}
.btn-text-green:hover {
    background: #d1e7dd;
    transform: translateY(-2px);
}

.stats-grid {
    display: grid;
    gap: 1.5rem;
    grid-template-columns: 1.4fr 0.8fr 0.8fr;
    margin-bottom: 2rem;
}
@media (max-width: 920px) {
    .stats-grid {
        grid-template-columns: 1fr;
    }
}

.card,
.section-container {
    background: white;
    padding: 1.5rem;
    border-radius: 12px;
    border: 1px solid #eee;
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.02);
    margin-bottom: 2rem;
}
.card--score {
    background: #2d7a4f;
    color: white;
    border: none;
}
.tag-score {
    background: rgba(255, 255, 255, 0.2);
    padding: 4px 10px;
    border-radius: 6px;
    font-size: 0.75rem;
    font-weight: 800;
    display: inline-block;
    margin-bottom: 1rem;
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
    display: block;
}
.card-num {
    font-size: 2.5rem;
    font-weight: bold;
    color: #1a1a1a;
    margin-bottom: 0.5rem;
    line-height: 1;
}
.text-dm {
    color: #666;
    margin-bottom: 1rem;
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
.status-logistique {
    background: #f0f0f0;
    color: #555;
    padding: 4px 10px;
    border-radius: 6px;
    font-size: 0.75rem;
    font-weight: bold;
    text-transform: uppercase;
}
.status-valid {
    background: #e9f5ed;
    color: #1e5636;
    padding: 4px 8px;
    border-radius: 4px;
    font-size: 0.8rem;
    font-weight: bold;
}
.status-neutral {
    background: #f5f5f5;
    color: #666;
    padding: 4px 8px;
    border-radius: 4px;
    font-size: 0.8rem;
    font-weight: bold;
}

.section-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1.5rem;
}
.data-table {
    width: 100%;
    border-collapse: collapse;
}
.data-table th {
    padding: 12px;
    text-align: left;
    border-bottom: 2px solid #eee;
    color: #666;
    font-weight: 600;
}
.data-table td {
    padding: 16px 12px;
    border-bottom: 1px solid #eee;
}
.table-subtext {
    color: #666;
    font-size: 0.8rem;
    display: block;
    margin-top: 4px;
}

.text-truncate {
    display: inline-block;
    max-width: 250px;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    color: #6d7b72;
}
.material-tag {
    background: #e9f5ed;
    color: #1e5636;
    padding: 4px 8px;
    border-radius: 6px;
    font-size: 0.85rem;
    font-weight: 700;
}
.stats-text {
    font-size: 0.85rem;
    color: #555;
}

.actions-cell {
    display: flex;
    gap: 8px;
    align-items: center;
}
.text-right {
    text-align: right !important;
}
.text-center {
    text-align: center;
}
.text-grey {
    color: #999;
}
.py-4 {
    padding-top: 1rem;
    padding-bottom: 1rem;
}

.planning-week {
    display: grid;
    grid-template-columns: repeat(7, 1fr);
    gap: 1rem;
}
.planning-day {
    background: #fbfdfb;
    border: 1px solid #e0e0e0;
    border-radius: 12px;
    padding: 10px;
    min-height: 100px;
}
.planning-day.is-today {
    border: 2px solid #2d7a4f;
    background: #eaf4ed;
}
.day-label {
    text-align: center;
    color: #666;
    font-size: 0.85rem;
    margin-bottom: 8px;
    display: flex;
    flex-direction: column;
}
.day-label strong {
    color: #1a1a1a;
    font-size: 1.1rem;
    margin-top: 4px;
}
.entry {
    padding: 6px;
    border-radius: 6px;
    font-size: 0.75rem;
    margin-bottom: 4px;
    display: flex;
    flex-direction: column;
}
.entry--formation {
    background: #d1e7dd;
    color: #0f5132;
}
.entry--event {
    background: #fff3cd;
    color: #856404;
}

.end-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
    gap: 1.5rem;
}
.mt-2 {
    margin-top: 0.5rem;
}
.mt-4 {
    margin-top: 1rem;
}
.mb-2 {
    margin-bottom: 0.5rem;
}
.state-card {
    border: 1px dashed #cfe0d4;
    border-radius: 14px;
    padding: 26px;
    text-align: center;
    color: #666;
}
</style>
