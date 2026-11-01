<script lang="ts">
  import type { ThemeId } from "../themes/theme";

  interface Props {
    theme: ThemeId;
    onThemeChange: (theme: ThemeId) => void;
  }

  let { theme, onThemeChange }: Props = $props();

  const themes: {
    id: ThemeId;
    name: string;
    description: string;
  }[] = [
    {
      id: "zen",
      name: "Zen",
      description: "Natural y tranquilo",
    },
    {
      id: "autumn",
      name: "Autumn",
      description: "Cálido y acogedor",
    },
    {
      id: "forest",
      name: "Forest",
      description: "Profundo y sereno",
    },
    {
      id: "midnight",
      name: "Midnight",
      description: "Nocturno y sobrio",
    },
    {
      id: "techno",
      name: "Techno",
      description: "Neón y cyberpunk",
    },

    {
      id: "ocean",
      name: "Ocean",
      description: "Azul fresco y luminoso",
    },
    {
      id: "lavender",
      name: "Lavender",
      description: "Violeta suave y sereno",
    },
    {
      id: "graphite",
      name: "Graphite",
      description: "Oscuro con acentos cálidos",
    },
  ];
</script>

<main class="content settings-page">
  <header class="page-header">
    <h1>Configuración</h1>
    <p>Personalizá la apariencia de Luma según tus preferencias.</p>
  </header>

  <section class="panel">
    <h2>Apariencia</h2>

    <p class="settings-description">
      Elegí el ambiente que más te guste. Podés cambiarlo cuando quieras.
    </p>

    <div class="theme-grid">
      {#each themes as option (option.id)}
        <button
          type="button"
          class="theme-option"
          class:selected={theme === option.id}
          aria-pressed={theme === option.id}
          onclick={() => onThemeChange(option.id)}
        >
          <span class="theme-preview {option.id}" aria-hidden="true"></span>

          <span class="theme-name">
            {option.name}
          </span>

          <span class="theme-description">
            {option.description}
          </span>
        </button>
      {/each}
    </div>
  </section>
</main>
