
<section class="welcome-banner">
  <div class="welcome-copy">
    <p class="welcome-kicker">Bienvenido a Luma</p>

    <h1>Un espacio para avanzar a tu ritmo.</h1>

    <p class="welcome-description">
      Organizá tus tareas y encontrá tu momento
      de concentración.
    </p>
  </div>

  <svg
    class="landscape"
    viewBox="0 0 640 240"
    preserveAspectRatio="xMidYMax slice"
    aria-hidden="true"
    focusable="false"
  >
    <!-- Sol -->
    <circle
      class="sun"
      cx="515"
      cy="55"
      r="25"
    />

    <!-- Nubes -->
    <g class="clouds">
      <path
        d="M390 69
           Q397 57 413 61
           Q424 43 444 55
           Q456 56 460 69 Z"
      />

      <path
        d="M545 98
           Q553 86 566 90
           Q578 75 596 89
           Q610 88 617 98 Z"
      />
    </g>

    <!-- Montañas lejanas -->
    <path
      class="mountains-back"
      d="M0 191
         L115 111
         Q125 104 135 113
         L213 175
         L314 83
         Q326 73 338 87
         L448 185
         L640 192
         V240 H0 Z"
    />

    <!-- Montañas intermedias -->
    <path
      class="mountains-middle"
      d="M0 209
         Q83 168 166 192
         Q253 143 332 190
         Q424 151 507 181
         Q581 165 640 186
         V240 H0 Z"
    />

    <!-- Valle -->
    <path
      class="valley"
      d="M0 221
         Q112 193 220 214
         Q345 181 457 208
         Q557 189 640 202
         V240 H0 Z"
    />

    <!-- Lago -->
    <path
      class="lake"
      d="M165 220
         Q287 202 391 216
         Q506 230 640 209
         V240 H165 Z"
    />

    <!-- Reflejos del agua -->
    <g class="water-lines">
      <path d="M300 221 H371" />
      <path d="M405 228 H485" />
      <path d="M520 218 H571" />
    </g>

    <!-- Árboles -->
    <g class="trees">
      <path d="M91 224 V191" />
      <path d="M77 202 L91 172 L105 202 Z" />

      <path d="M133 227 V202" />
      <path d="M122 211 L133 184 L144 211 Z" />

      <path d="M480 222 V185" />
      <path d="M464 199 L480 165 L496 199 Z" />

      <path d="M522 225 V192" />
      <path d="M508 205 L522 174 L536 205 Z" />

      <path d="M588 224 V178" />
      <path d="M570 194 L588 154 L606 194 Z" />
    </g>

    <!-- Primer plano -->
    <path
      class="foreground"
      d="M0 235
         Q85 216 166 231
         Q232 219 295 236
         Q385 226 452 235
         Q550 218 640 230
         V240 H0 Z"
    />
  </svg>
</section>

<style>
  .welcome-banner {
    position: relative;
    display: flex;
    align-items: center;

    min-height: 235px;
    margin-bottom: 28px;
    padding: 36px 40px;

    overflow: hidden;
    isolation: isolate;

    background: var(--color-surface-soft);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-card);
  }

  /* Texto principal */

  .welcome-copy {
    position: relative;
    z-index: 2;
    width: 54%;
    max-width: 490px;
  }

  .welcome-kicker {
    margin: 0 0 12px;

    color: var(--color-primary);
    font-size: 12px;
    font-weight: 700;
    letter-spacing: 1.2px;
    text-transform: uppercase;
  }

  h1 {
    margin: 0 0 12px;

    color: var(--color-text);
    font-size: clamp(25px, 3vw, 35px);
    line-height: 1.2;
    letter-spacing: -0.8px;
  }

  .welcome-description {
    max-width: 380px;
    margin: 0;

    color: var(--color-text-muted);
    font-size: 14px;
    line-height: 1.6;
  }

  /* Paisaje SVG */

  .landscape {
    position: absolute;
    right: 0;
    bottom: 0;

    width: 58%;
    height: 100%;

    z-index: 1;
    pointer-events: none;

    -webkit-mask-image: linear-gradient(
      to right,
      transparent,
      black 23%
    );
    mask-image: linear-gradient(
      to right,
      transparent,
      black 23%
    );
  }

  /* Colores derivados del tema */

  .sun {
    fill: #F1DB9D;
    opacity: 0.75;
  }

  .clouds {
    fill: var(--color-surface);
    opacity: 0.55;
  }

  .mountains-back {
    fill: color-mix(
      in srgb,
      var(--color-primary) 22%,
      var(--color-surface-soft)
    );
  }

  .mountains-middle {
    fill: color-mix(
      in srgb,
      var(--color-primary) 42%,
      var(--color-surface-soft)
    );
  }

  .valley {
    fill: color-mix(
      in srgb,
      var(--color-primary) 60%,
      var(--color-surface-soft)
    );
  }

  .lake {
    fill: color-mix(
      in srgb,
      var(--color-primary) 20%,
      var(--color-surface)
    );
  }

  .water-lines {
    fill: none;
    stroke: var(--color-surface);
    stroke-width: 2;
    stroke-linecap: round;
    opacity: 0.6;
  }

  .trees {
    fill: color-mix(
      in srgb,
      var(--color-primary) 85%,
      var(--color-background)
    );
    stroke: color-mix(
      in srgb,
      var(--color-primary) 85%,
      var(--color-background)
    );
    stroke-width: 3;
    stroke-linecap: round;
    stroke-linejoin: round;
  }

  .foreground {
    fill: color-mix(
      in srgb,
      var(--color-primary) 70%,
      var(--color-surface-soft)
    );
  }

  /* Adaptación a dispositivos móviles */

  @media (max-width: 700px) {
    .welcome-banner {
      min-height: 330px;
      padding: 26px;
      align-items: flex-start;
    }

    .welcome-copy {
      width: 100%;
      max-width: 100%;
    }

    h1 {
      font-size: clamp(24px, 6vw, 30px);
    }

    .landscape {
      width: 100%;
      height: 155px;

      -webkit-mask-image: linear-gradient(
        to bottom,
        transparent,
        black 35%
      );
      mask-image: linear-gradient(
        to bottom,
        transparent,
        black 35%
      );
    }
  }
</style>
