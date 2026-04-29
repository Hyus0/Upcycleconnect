<template>
    <header class="content-header">
        <div class="header-left">
            <p class="sidebar__category2">ACCUEIL > TABLEAU DE BORD</p>
            <h1 class="hero-title1">Bonjour {{ prenom }} 👋</h1>
            <p class="classic-text">
                Voici un resume de votre activite sur UpcycleConnect
            </p>
        </div>
        <router-link to="/profil/createAnnonce" class="btn-main-action">+ Deposer une annonce</router-link>
    </header>

    <div class="stats-grid">
        <div class="card card--score">
            <p class="tag-score">UPCYCLING SCORE</p>

            <div class="score-value">
                {{ stats.total_points }} <span>pts</span>
            </div>

            <p class="score-level">Niveau : {{ stats.niveau }}</p>

            <div class="score-footer">
                <div class="mini-stat">
                    <strong>{{ stats.co2_total_evite_kg }} kg</strong><br />
                    CO2 evite
                </div>

                <div class="mini-stat">
                    <strong>{{ stats.nb_objets_recycles }}</strong><br />
                    Objets recycles
                </div>

                <div class="mini-stat">
                    <strong>EUR {{ stats.ressources_economisees }}</strong><br />
                    Economise
                </div>
            </div>
        </div>
        <div class="card card--white">
            <div class="card-num">{{ annoncesActivesCount }}</div>
            <p class="text-dm">Annonces actives</p>
            <span class="badge badge--green">{{ annoncesActivesCount > 0 ? "+1 ce mois" : "Aucune active" }}</span>
        </div>

        <div class="card card--white">
            <div class="card-num2">{{ calendarEntries.length }}</div>
            <p class="text-dm">Inscriptions a venir</p>
            <span class="badge badge--orange">{{ nextEntryLabel }}</span>
        </div>
    </div>

    <div class="section-container">
        <div class="section-header">
            <h2>Mes dernieres annonces</h2>
            <div class="header-actions">
                <input
                    type="text"
                    placeholder="Rechercher..."
                    class="search-input"
                />
                <button class="btn-secondary">Tous statuts</button>
                <router-link to="/profil/createAnnonce" class="btn-main-action1">+ Nouvelle annonce</router-link>
            </div>
        </div>
        <table class="data-table">
            <thead>
                <tr>
                    <th>OBJET</th>
                    <th>CATEGORIE</th>
                    <th>TYPE</th>
                    <th>STATUT</th>
                    <th>DATE</th>
                    <th>ACTIONS</th>
                </tr>
            </thead>
            <tbody>
                <tr v-for="annonce in annonces.slice(0, 4)" :key="annonce.id">
                    <td>
                        <strong>{{ annonce.titre }}</strong><br>
                        <small>{{ annonce.type_materiau }}</small>
                    </td>

                    <td>
                        <span :class="annonce.type === 'Don' ? 'tag-don' : 'tag-vente'">
                            {{ annonce.type === 'Don' ? 'DON' : 'VENTE ' + annonce.prix + 'EUR' }}
                        </span>
                    </td>

                    <td>
                        <span :class="annonce.est_valide === 'Valide' ? 'status-valid' : 'status-pending'">
                            {{ annonce.est_valide === 'Valide' ? 'APPROUVEE' : 'EN ANALYSE' }}
                        </span>
                    </td>

                    <td>
                        <span class="status-neutral">
                            {{ annonce.statut }}
                        </span>
                    </td>

                    <td>{{ formatDate(annonce.date_creation) }}</td>

                    <td class="actions-cell">
                        <button
                            v-if="annonce.est_valide === 'Valide' && annonce.statut === 'Disponible'"
                            class="btn-plan"
                            @click="goToPlanning(annonce.id)"
                        >
                            Planifier depot
                        </button>

                        <button
                            class="btn-view"
                            @click="goToAnnonce(annonce.id)"
                        >
                            Voir
                        </button>

                        <button
                            v-if="annonce.est_valide === 'En attente'"
                            class="btn-modify"
                            @click="goToModify(annonce.id)"
                        >
                            Modifier
                        </button>

                        <button class="btn-remove" @click="removeAnnonce(annonce.id)">Retirer</button>
                    </td>
                </tr>
            </tbody>
        </table>
        <p v-if="annonces.length === 0" class="empty-msg">Vous n'avez pas encore depose d'annonces.</p>
    </div>

    <div class="section-container planning-section">
        <div class="section-header planning-section__header">
            <div>
                <h2>Mon calendrier d'inscriptions</h2>
                <p class="planning-subtitle">
                    Formations, evenements et autres rendez-vous auxquels vous etes inscrit.
                </p>
            </div>
            <div class="planning-toolbar">
                <button class="btn-secondary" @click="previousMonth">Mois precedent</button>
                <span class="planning-month">{{ currentMonthLabel }}</span>
                <button class="btn-secondary" @click="nextMonth">Mois suivant</button>
            </div>
        </div>

        <div class="planning-legend">
            <span class="planning-legend__item"><i class="planning-dot planning-dot--formation"></i> Formation</span>
            <span class="planning-legend__item"><i class="planning-dot planning-dot--event"></i> Evenement</span>
            <span class="planning-legend__item"><i class="planning-dot planning-dot--other"></i> Autre rendez-vous</span>
        </div>

        <div class="planning-calendar">
            <div class="planning-calendar__weekday" v-for="day in weekDays" :key="day">{{ day }}</div>
            <div
                v-for="day in calendarDays"
                :key="day.key"
                class="planning-day"
                :class="{ 'is-outside': !day.isCurrentMonth, 'is-today': day.isToday, 'has-events': day.entries.length > 0 }"
            >
                <div class="planning-day__head">
                    <span class="planning-day__number">{{ day.label }}</span>
                    <small v-if="day.entries.length">{{ day.entries.length }} inscrit.</small>
                </div>

                <div class="planning-day__entries">
                    <div
                        v-for="entry in day.entries.slice(0, 3)"
                        :key="entry.id"
                        class="planning-entry"
                        :class="`planning-entry--${entry.kind}`"
                    >
                        <strong>{{ entry.title }}</strong>
                        <span>{{ entry.timeLabel }}</span>
                    </div>
                    <div v-if="day.entries.length > 3" class="planning-entry planning-entry--more">
                        +{{ day.entries.length - 3 }} autre(s)
                    </div>
                </div>
            </div>
        </div>

        <div v-if="calendarEntries.length === 0" class="planning-empty">
            <p>Aucune inscription n'est encore visible pour ce compte.</p>
            <span>
                Le calendrier se remplira des que vous rejoindrez une formation, un evenement ou un autre rendez-vous.
            </span>
        </div>
    </div>

    <div class="end-grid">
        <div class="section-container-tips">
            <p class="tag-vente">Conseil du jour</p>

            <div class="section-header">
                <h2>Transformer un vieux jean en sac</h2>
            </div>
            <p>
                Apprenez a confectionner un sac tote en 30 minutes avec un jean use.
                Materiel necessaire : aiguille, fil, ciseaux.
            </p>
            <button class="btn-view">Lire la suite -></button>
        </div>
        <div class="section-container-tips">
            <p class="tag-don">Notification</p>

            <div class="section-header">
                <h2>Votre depot a ete recupere !</h2>
            </div>
            <p>
                La chaise vintage que vous avez deposee le 12 fev. a ete recuperee par un artisan.
                +50 points Upcycling Score !
            </p>
            <button class="btn-view">Voir le projet -></button>
        </div>
    </div>
</template>

<script setup>
import { computed, onMounted, ref } from "vue";
import { useRouter } from "vue-router";

const router = useRouter();

const prenom = ref(localStorage.getItem("userPrenom") || "Invite");
const annonces = ref([]);
const calendarEntries = ref([]);
const currentMonth = ref(new Date(new Date().getFullYear(), new Date().getMonth(), 1));
const weekDays = ["Lun", "Mar", "Mer", "Jeu", "Ven", "Sam", "Dim"];

const stats = ref({
    total_points: 0,
    niveau: "Chargement...",
    co2_total_evite_kg: 0,
    nb_objets_recycles: 0,
    ressources_economisees: 0,
});

const annoncesActivesCount = computed(() =>
    annonces.value.filter((annonce) => annonce.est_valide === "Valide" && annonce.statut === "Disponible").length
);

const nextEntryLabel = computed(() => {
    if (calendarEntries.value.length === 0) return "A venir";
    const nextEntry = [...calendarEntries.value].sort((a, b) => a.date.localeCompare(b.date))[0];
    return nextEntry.dateLabel;
});

const currentMonthLabel = computed(() =>
    currentMonth.value.toLocaleDateString("fr-FR", { month: "long", year: "numeric" })
);

const calendarDays = computed(() => {
    const monthStart = new Date(currentMonth.value);
    const start = new Date(monthStart);
    const offset = (monthStart.getDay() + 6) % 7;
    start.setDate(monthStart.getDate() - offset);

    return Array.from({ length: 42 }, (_, index) => {
        const date = new Date(start);
        date.setDate(start.getDate() + index);
        const isoDate = date.toISOString().slice(0, 10);

        return {
            key: isoDate,
            label: date.getDate(),
            isoDate,
            isCurrentMonth: date.getMonth() === currentMonth.value.getMonth(),
            isToday: isoDate === new Date().toISOString().slice(0, 10),
            entries: calendarEntries.value.filter((entry) => entry.date === isoDate)
        };
    });
});

const formatDate = (dateString) => {
    if (!dateString) return "...";
    const date = new Date(dateString);
    return date.toLocaleDateString("fr-FR", {
        day: "numeric",
        month: "short",
        year: "numeric",
    });
};

const parseDate = (value) => {
    if (!value) return null;
    const parsed = new Date(value);
    if (Number.isNaN(parsed.getTime())) return null;
    return parsed;
};

const toCalendarEntry = (item, kind, title, dateValue) => {
    const parsed = parseDate(dateValue);
    if (!parsed) return null;

    return {
        id: `${kind}-${item.id ?? title}-${parsed.getTime()}`,
        kind,
        title,
        date: parsed.toISOString().slice(0, 10),
        dateLabel: parsed.toLocaleDateString("fr-FR", { day: "numeric", month: "short" }),
        timeLabel: parsed.toLocaleTimeString("fr-FR", { hour: "2-digit", minute: "2-digit" })
    };
};

const readJoinedIds = (key) => {
    try {
        const raw = localStorage.getItem(key);
        const parsed = JSON.parse(raw || "[]");
        return Array.isArray(parsed) ? parsed.map((value) => Number(value)) : [];
    } catch {
        return [];
    }
};

const readCustomCalendarEntries = (key) => {
    try {
        const raw = localStorage.getItem(key);
        const parsed = JSON.parse(raw || "[]");
        if (!Array.isArray(parsed)) return [];

        return parsed
            .map((entry, index) => {
                const parsedDate = parseDate(entry.date || entry.date_debut || entry.datetime);
                if (!parsedDate) return null;

                return {
                    id: `other-${entry.id ?? index}`,
                    kind: "other",
                    title: entry.title ?? entry.titre ?? "Rendez-vous",
                    date: parsedDate.toISOString().slice(0, 10),
                    dateLabel: parsedDate.toLocaleDateString("fr-FR", { day: "numeric", month: "short" }),
                    timeLabel: parsedDate.toLocaleTimeString("fr-FR", { hour: "2-digit", minute: "2-digit" })
                };
            })
            .filter(Boolean);
    } catch {
        return [];
    }
};

const goToAnnonce = (id) => {
    router.push({ name: "see-annonce", params: { id } });
};

const goToModify = (id) => {
    router.push({ name: "modification-annonce", params: { id } });
};

const goToPlanning = (id) => {
    router.push({ name: "mes-depots", query: { selectedAnnonce: id } });
};

const removeAnnonce = (id) => {
    console.log("Demande de suppression pour :", id);
};

const previousMonth = () => {
    currentMonth.value = new Date(currentMonth.value.getFullYear(), currentMonth.value.getMonth() - 1, 1);
};

const nextMonth = () => {
    currentMonth.value = new Date(currentMonth.value.getFullYear(), currentMonth.value.getMonth() + 1, 1);
};

const loadCalendarEntries = async (id, token) => {
    const joinedFormationIds = readJoinedIds(`upcycleconnect-joined-formations-${id}`);
    const joinedEventIds = readJoinedIds(`upcycleconnect-joined-events-${id}`);
    const customEntries = readCustomCalendarEntries(`upcycleconnect-calendar-${id}`);

    const [formationsResponse, eventsResponse] = await Promise.all([
        fetch("http://localhost:8081/formations", { headers: { Authorization: token } }).catch(() => null),
        fetch("http://localhost:8081/evenements", { headers: { Authorization: token } }).catch(() => null)
    ]);

    const formationItems = formationsResponse?.ok ? await formationsResponse.json() : [];
    const eventItems = eventsResponse?.ok ? await eventsResponse.json() : [];

    const formationEntries = Array.isArray(formationItems)
        ? formationItems
              .filter((item) => joinedFormationIds.includes(Number(item.id)))
              .map((item) => toCalendarEntry(item, "formation", item.titre ?? "Formation", item.date_debut))
              .filter(Boolean)
        : [];

    const eventEntries = Array.isArray(eventItems)
        ? eventItems
              .filter((item) => joinedEventIds.includes(Number(item.id)))
              .map((item) => toCalendarEntry(item, "event", item.titre ?? "Evenement", item.date_evenement ?? item.date))
              .filter(Boolean)
        : [];

    calendarEntries.value = [...formationEntries, ...eventEntries, ...customEntries].sort((a, b) =>
        a.date.localeCompare(b.date)
    );
};

onMounted(async () => {
    const id = localStorage.getItem("userId");
    const token = localStorage.getItem("userToken");

    if (!id || !token) return;

    try {
        const response = await fetch(`http://localhost:8081/users/${id}/stats`, {
            headers: { Authorization: token },
        });
        if (response.ok) {
            const data = await response.json();
            stats.value = data;
            if (data.total_points !== undefined) {
                localStorage.setItem("userScore", data.total_points.toString());
                window.dispatchEvent(new Event("auth-change"));
            }
        }
    } catch (error) {
        console.error("Erreur stats :", error);
    }

    try {
        const resAnnonces = await fetch(`http://localhost:8081/users/${id}/annonces`, {
            headers: { Authorization: token }
        });
        if (resAnnonces.ok) {
            annonces.value = await resAnnonces.json();
        }
    } catch (error) {
        console.error("Erreur annonces :", error);
    }

    try {
        await loadCalendarEntries(id, token);
    } catch (error) {
        console.error("Erreur calendrier :", error);
        calendarEntries.value = [];
    }
});
</script>
