<script lang="ts">
  import Sidebar from "./lib/components/Sidebar.svelte";
  import HomePage from "./lib/pages/HomePage.svelte";
  import TasksPage from "./lib/pages/TasksPage.svelte";
  import CalendarPage from "./lib/pages/CalendarPage.svelte";
  import SearchPage from "./lib/pages/SearchPage.svelte";
  import TrashPage from "./lib/pages/TrashPage.svelte";
  import SettingsPage from "./lib/pages/SettingsPage.svelte";

  import type { PageId } from "./lib/navigation";
  import type { ThemeId } from "./lib/themes/theme";
  import { loadTheme } from "./lib/themes/theme";

  let currentPage = $state<PageId>("home");
  let theme = $state<ThemeId>(loadTheme());

  function navigateTo(page: PageId) {
    currentPage = page;
  }

  function changeTheme(newTheme: ThemeId) {
    theme = newTheme;
  }

  $effect(() => {
    document.documentElement.dataset.theme = theme;

    try {
      localStorage.setItem("luma.theme", theme);
    } catch {
      // La aplicación sigue funcionando sin almacenamiento.
    }
  });
</script>

<svelte:head>
  <title>Luma</title>
  <meta
    name="description"
    content="Luma, un espacio tranquilo para organizar tus tareas."
  />
</svelte:head>

<div class="layout">
  <Sidebar {currentPage} onNavigate={navigateTo} />

  <div class="app-main">
    {#if currentPage === "home"}
      <HomePage />
    {:else if currentPage === "tasks"}
      <TasksPage />
    {:else if currentPage === "calendar"}
      <CalendarPage />
    {:else if currentPage === "search"}
      <SearchPage />
    {:else if currentPage === "trash"}
      <TrashPage />
    {:else if currentPage === "settings"}
      <SettingsPage {theme} onThemeChange={changeTheme} />
    {/if}

    <footer class="app-footer">
      <div class="app-footer-inner">
        <span>© 2026 Luma</span>
        <span>Organización y gestión de tareas</span>
      </div>
    </footer>
  </div>
</div>
