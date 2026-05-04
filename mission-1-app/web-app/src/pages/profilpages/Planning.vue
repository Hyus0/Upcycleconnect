<template>
    <div class="layout-wrapper planning-page">
        <header class="content-header">
            <div class="header-left">
                <p class="sidebar__category2">ACCUEIL > MON PLANNING</p>
                <h1 class="hero-title1">Calendrier d'inscriptions</h1>
            </div>
            <router-link to="/profil" class="btn-secondary" style="text-decoration: none"
                >🠔 Retour au tableau de bord</router-link
            >
        </header>

        <div class="section-container">
            <div class="planning-toolbar">
                <div class="toolbar-left">
                    <button class="btn-icon" @click="previousMonth">
                        ❮ Précédent
                    </button>
                    <button class="btn-today" @click="goToToday">
                        Aujourd'hui
                    </button>
                    <button class="btn-icon" @click="nextMonth">
                        Suivant ❯
                    </button>
                </div>
                <h2 class="planning-month">{{ currentMonthLabel }}</h2>
                <div class="planning-legend">
                    <span class="planning-legend__item"
                        ><i class="planning-dot planning-dot--formation"></i>
                        Formation</span
                    >
                    <span class="planning-legend__item"
                        ><i class="planning-dot planning-dot--event"></i>
                        Événement</span
                    >
                </div>
            </div>

            <div class="planning-calendar">
                <div
                    class="planning-calendar__weekday"
                    v-for="day in weekDays"
                    :key="day"
                >
                    {{ day }}
                </div>

                <div
                    v-for="day in calendarDays"
                    :key="day.key"
                    class="planning-day"
                    :class="{
                        'is-outside': !day.isCurrentMonth,
                        'is-today': day.isToday,
                    }"
                >
                    <div class="planning-day__head">
                        <span class="planning-day__number">{{
                            day.label
                        }}</span>
                    </div>

                    <div class="planning-day__entries">
                        <div
                            v-for="entry in day.entries"
                            :key="entry.id"
                            class="planning-entry clickable-entry"
                            :class="`planning-entry--${entry.kind}`"
                            @click="openEntryDetail(entry)"
                        >
                            <strong>{{ entry.title }}</strong>
                            <span>{{ entry.timeLabel }}</span>
                        </div>
                    </div>
                </div>
            </div>
        </div>

        <!-- MODAL DÉTAIL D'UNE ENTRÉE (Pop-up) -->
        <div
            v-if="selectedEntry"
            class="planning-modal"
            @click.self="closeEntryDetail"
        >
            <div class="detail-panel">
                <button class="btn-close-top" @click="closeEntryDetail">
                    ✖
                </button>
                <div class="detail-header">
                    <span class="tag-type" :class="`tag-${selectedEntry.kind}`">
                        {{
                            selectedEntry.kind === "formation"
                                ? "FORMATION"
                                : "ÉVÉNEMENT"
                        }}
                    </span>
                    <h2 class="detail-title">{{ selectedEntry.title }}</h2>
                </div>

                <div class="detail-body">
                    <div class="detail-row">
                        <span class="icon">📅</span>
                        <div>
                            <strong>Date :</strong>
                            <p>{{ selectedEntry.dateLabelLong }}</p>
                        </div>
                    </div>
                    <div class="detail-row">
                        <span class="icon">⏰</span>
                        <div>
                            <strong>Heure :</strong>
                            <p>{{ selectedEntry.timeLabel }}</p>
                        </div>
                    </div>
                </div>

                <button
                    class="btn-main-action-full"
                    @click="goToEntryRoute(selectedEntry)"
                >
                    Voir la page de l'{{
                        selectedEntry.kind === "formation"
                            ? "formation"
                            : "événement"
                    }}
                </button>
            </div>
        </div>
    </div>
</template>

<script setup>
import { computed, onMounted, ref } from "vue";
import { useRouter } from "vue-router";

const router = useRouter();
const API_URL = "http://localhost:8081";

const calendarEntries = ref([]);
const currentMonth = ref(
    new Date(new Date().getFullYear(), new Date().getMonth(), 1),
);
const weekDays = ["Lun", "Mar", "Mer", "Jeu", "Ven", "Sam", "Dim"];
const selectedEntry = ref(null);

const getLocalISODate = (date) => {
    const d = new Date(date);
    const year = d.getFullYear();
    const month = String(d.getMonth() + 1).padStart(2, "0");
    const day = String(d.getDate()).padStart(2, "0");
    return `${year}-${month}-${day}`;
};

const currentMonthLabel = computed(() => {
    const label = currentMonth.value.toLocaleDateString("fr-FR", {
        month: "long",
        year: "numeric",
    });
    return label.charAt(0).toUpperCase() + label.slice(1);
});

const calendarDays = computed(() => {
    const year = currentMonth.value.getFullYear();
    const month = currentMonth.value.getMonth();

    const firstDayOfMonth = new Date(year, month, 1);
    let offset = firstDayOfMonth.getDay() - 1;
    if (offset === -1) offset = 6;

    const start = new Date(year, month, 1 - offset);

    return Array.from({ length: 42 }, (_, index) => {
        const date = new Date(
            start.getFullYear(),
            start.getMonth(),
            start.getDate() + index,
        );
        const isoDate = getLocalISODate(date);

        return {
            key: isoDate,
            label: date.getDate(),
            isoDate,
            isCurrentMonth: date.getMonth() === month,
            isToday: isoDate === getLocalISODate(new Date()),
            entries: calendarEntries.value.filter(
                (entry) => entry.date === isoDate,
            ),
        };
    });
});

const previousMonth = () =>
    (currentMonth.value = new Date(
        currentMonth.value.getFullYear(),
        currentMonth.value.getMonth() - 1,
        1,
    ));
const nextMonth = () =>
    (currentMonth.value = new Date(
        currentMonth.value.getFullYear(),
        currentMonth.value.getMonth() + 1,
        1,
    ));
const goToToday = () =>
    (currentMonth.value = new Date(
        new Date().getFullYear(),
        new Date().getMonth(),
        1,
    ));

const openEntryDetail = (entry) => (selectedEntry.value = entry);
const closeEntryDetail = () => (selectedEntry.value = null);

const goToEntryRoute = (entry) => {
    closeEntryDetail();
    if (entry.kind === "formation") {
        router.push({
            name: "FormationDetail",
            params: { id: entry.originalId },
        });
    } else if (entry.kind === "event") {
        router.push({
            name: "EvenementDetail",
            params: { id: entry.originalId },
        });
    }
};

const toCalendarEntry = (item, kind, title, dateValue) => {
    if (!dateValue) return null;
    const parsed = new Date(dateValue);
    if (Number.isNaN(parsed.getTime())) return null;

    return {
        id: `${kind}-${item.id ?? title}-${parsed.getTime()}`,
        originalId: item.id,
        kind,
        title,
        date: getLocalISODate(parsed),
        dateLabelLong: parsed.toLocaleDateString("fr-FR", {
            weekday: "long",
            day: "numeric",
            month: "long",
            year: "numeric",
        }),
        timeLabel: parsed.toLocaleTimeString("fr-FR", {
            hour: "2-digit",
            minute: "2-digit",
        }),
    };
};

const loadCalendarEntries = async (id) => {
    try {
        const res = await fetch(`${API_URL}/user/planning/${id}`, {
            headers: {
                Authorization: `Bearer ${localStorage.getItem("userToken") || ""}`,
            },
        });

        if (!res.ok) return;

        const data = await res.json();
        let entries = [];

        if (data && !Array.isArray(data)) {
            const formations = (data.formations || []).map((f) =>
                toCalendarEntry(f, "formation", f.titre, f.date_debut),
            );
            const evenements = (data.evenements || []).map((e) =>
                toCalendarEntry(e, "event", e.titre, e.date_evenement),
            );
            entries = [...formations, ...evenements].filter(Boolean);
        }

        calendarEntries.value = entries.sort((a, b) =>
            a.date.localeCompare(b.date),
        );
    } catch (error) {
        console.error("Erreur calendrier :", error);
    }
};

onMounted(() => {
    const id = localStorage.getItem("userId");
    if (id) loadCalendarEntries(id);
});
</script>

<style scoped>
.planning-page {
    padding: 20px;
    background: #f7f9f7;
    min-height: 100vh;
}

.content-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 2rem;
}


.section-container {
    background: white;
    padding: 2rem;
    border-radius: 16px;
    border: 1px solid #eee;
}

.planning-toolbar {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 2rem;
}

.toolbar-left {
    display: flex;
    gap: 0.5rem;
}

.btn-icon,
.btn-today {
    padding: 8px 16px;
    border-radius: 8px;
    border: 1px solid #ddd;
    background: white;
    cursor: pointer;
    font-weight: 600;
}

.planning-month {
    font-size: 1.5rem;
    font-weight: 900;
    margin: 0;
}

.planning-calendar {
    display: grid;
    grid-template-columns: repeat(7, 1fr);
    gap: 10px;
}

.planning-calendar__weekday {
    text-align: center;
    font-weight: bold;
    color: #888;
    padding-bottom: 10px;
}

.planning-day {
    border: 1px solid #eaeaea;
    border-radius: 8px;
    min-height: 120px;
    padding: 10px;
}

.planning-day.is-outside {
    background: #fafafa;
    opacity: 0.6;
}

.planning-day.is-today {
    border: 2px solid #2d7a4f;
    background: #f4fbf7;
}

.planning-day__head {
    font-weight: bold;
    margin-bottom: 8px;
}

.planning-entry {
    padding: 6px;
    border-radius: 6px;
    font-size: 0.8rem;
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

.clickable-entry {
    cursor: pointer;
}

.clickable-entry:hover {
    filter: brightness(0.95);
}

.planning-modal {
    position: fixed;
    top: 0;
    left: 0;
    width: 100vw;
    height: 100vh;
    background: rgba(0, 0, 0, 0.5);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 1000;
}

.detail-panel {
    background: white;
    width: 400px;
    padding: 2rem;
    border-radius: 12px;
    position: relative;
}

.btn-close-top {
    position: absolute;
    top: 10px;
    right: 10px;
    background: none;
    border: none;
    font-size: 1.2rem;
    cursor: pointer;
}

.tag-type {
    font-size: 0.75rem;
    font-weight: bold;
    padding: 4px 8px;
    border-radius: 4px;
    display: inline-block;
    margin-bottom: 10px;
}
.tag-formation {
    background: #d1e7dd;
    color: #0f5132;
}
.tag-event {
    background: #fff3cd;
    color: #856404;
}

.detail-title {
    margin-bottom: 1rem;
    font-size: 1.3rem;
}
.detail-body {
    background: #f9f9f9;
    padding: 1rem;
    border-radius: 8px;
    margin-bottom: 1.5rem;
}
.detail-row {
    display: flex;
    align-items: center;
    gap: 10px;
    margin-bottom: 10px;
}
.detail-row p {
    margin: 0;
}

.planning-legend {
    display: flex;
    gap: 1.5rem;
}

.planning-legend__item {
    font-size: 0.9rem;
    font-weight: bold;
    color: #666;
    display: flex;
    align-items: center;
    gap: 8px;
}

.planning-dot {
    width: 12px;
    height: 12px;
    border-radius: 50%;
    display: inline-block;
}

.planning-dot--formation {
    background: #d1e7dd;
    border: 2px solid #0f5132;
}
.planning-dot--event {
    background: #fff3cd;
    border: 2px solid #856404;
}

.btn-main-action-full {
    width: 100%;
    padding: 12px;
    background: #2d7a4f;
    color: white;
    border: none;
    border-radius: 8px;
    font-weight: bold;
    cursor: pointer;
}
</style>
