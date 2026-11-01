<section class="panel focus-widget" aria-labelledby="focus-title">
  <header class="widget-header">
    <div>
      <p class="widget-kicker">Mi espacio · Herramienta 01</p>
      <h2 id="focus-title">Concentración</h2>
    </div>

    <span class="widget-badge">Pomodoro</span>
  </header>

  <div class="scenic-card">
    <div class="timer-content">
      <span class="session-type">Sesión libre</span>

      <strong class="timer-value">25:00</strong>

      <p>Un momento para enfocarte.</p>
    </div>

    <svg
      class="landscape"
      viewBox="0 0 480 260"
      preserveAspectRatio="xMidYMax slice"
      aria-hidden="true"
      focusable="false"
    >
      <!-- Sol -->
      <circle
        class="scene-sun"
        cx="390"
        cy="68"
        r="26"
      />

      <!-- Nubes -->
      <g class="clouds">
        <path
          d="M285 85
             Q293 73 307 77
             Q317 61 335 72
             Q348 71 353 85 Z"
        />
        <path
          d="M360 105
             Q367 95 379 99
             Q389 86 405 97
             Q418 97 423 105 Z"
        />
      </g>

      <!-- Montañas -->
      <path
        class="mountains-back"
        d="M0 190
           L92 112
           Q101 103 110 112
           L176 171
           L262 100
           Q273 90 284 101
           L370 176
           L480 184
           V260 H0 Z"
      />

      <path
        class="mountains-middle"
        d="M0 205
           Q80 170 152 192
           Q224 154 298 194
           Q380 157 480 188
           V260 H0 Z"
      />

      <!-- Lago -->
      <path
        class="lake"
        d="M148 214
           Q220 206 279 213
           Q346 222 415 214
           Q447 211 480 205
           V260 H148 Z"
      />

      <g class="water-lines">
        <path d="M230 222 H275" />
        <path d="M300 229 H355" />
        <path d="M380 220 H425" />
      </g>

      <!-- Colinas -->
      <path
        class="hill-back"
        d="M0 212
           Q78 181 161 206
           Q230 193 300 214
           Q382 196 480 210
           V260 H0 Z"
      />

      <path
        class="hill-front"
        d="M0 228
           Q68 202 138 222
           Q226 242 312 224
           Q392 208 480 223
           V260 H0 Z"
      />

      <!-- Árboles -->
      <g class="trees">
        <path d="M58 230 V193" />
        <path d="M42 207 L58 171 L74 207 Z" />

        <path d="M102 236 V205" />
        <path d="M88 216 L102 185 L116 216 Z" />

        <path d="M362 229 V188" />
        <path d="M345 202 L362 165 L379 202 Z" />

        <path d="M410 234 V199" />
        <path d="M394 213 L410 179 L426 213 Z" />
      </g>
    </svg>
  </div>

  <footer class="widget-footer">
    <p id="pomodoro-preview-note">
      Disponible en el Hito 3.
    </p>

    <button
      type="button"
      class="start-button"
      disabled
      aria-describedby="pomodoro-preview-note"
    >
      Iniciar sesión
    </button>
  </footer>
</section>

<style>
  .focus-widget {
    min-width: 0;
  }

  .widget-header {
    display: flex;
    align-items: flex-start;
    justify-content: space-between;
    gap: 12px;
    margin-bottom: 24px;
  }

  .widget-kicker {
    margin: 0 0 8px;
    color: var(--color-text-muted);
    font-size: 10px;
    font-weight: 600;
    letter-spacing: 0.8px;
    text-transform: uppercase;
  }

  .widget-header h2 {
    margin: 0;
  }

  .widget-badge {
    padding: 7px 12px;
    border-radius: 999px;
    background: var(--color-surface-soft);
    color: var(--color-primary);
    font-size: 12px;
    font-weight: 600;
  }

  .scenic-card {
    position: relative;
    min-height: 310px;
    overflow: hidden;
    border-radius: 18px;
    background: var(--color-surface-soft);
    isolation: isolate;
    border: 1px solid color-mix(
      in srgb,
      var(--color-border) 72%,
      transparent
    );
  }

  .timer-content {
    position: relative;
    z-index: 2;
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: 36px 16px 0;
    text-align: center;
  }

  .session-type {
    color: var(--color-primary);
    font-size: 13px;
    font-weight: 600;
  }

  .timer-value {
    margin: 12px 0 10px;
    color: var(--color-text);
    font-size: clamp(46px, 5vw, 64px);
    font-weight: 650;
    letter-spacing: -2px;
    font-variant-numeric: tabular-nums;
  }

  .timer-content p {
    margin: 0;
    color: var(--color-text-muted);
    font-size: 13px;
  }

  .landscape {
    position: absolute;
    inset: 0;
    width: 100%;
    height: 100%;
    z-index: 1;
    pointer-events: none;
  }

  .scene-sun {
    fill: #F1DB9D;
    opacity: 0.72;
  }

  .clouds {
    fill: var(--color-surface);
    opacity: 0.5;
  }

  .mountains-back {
    fill: color-mix(
      in srgb,
      var(--color-primary) 18%,
      var(--color-surface-soft)
    );
  }

  .mountains-middle {
    fill: color-mix(
      in srgb,
      var(--color-primary) 30%,
      var(--color-surface-soft)
    );
  }

  .lake {
    fill: color-mix(
      in srgb,
      var(--color-primary) 16%,
      var(--color-surface)
    );
    opacity: 0.95;
  }

  .water-lines {
    fill: none;
    stroke: var(--color-surface);
    stroke-width: 2;
    stroke-linecap: round;
    opacity: 0.55;
  }

  .hill-back {
    fill: color-mix(
      in srgb,
      var(--color-primary) 44%,
      var(--color-surface-soft)
    );
  }

  .hill-front {
    fill: color-mix(
      in srgb,
      var(--color-primary) 62%,
      var(--color-surface-soft)
    );
  }

  .trees {
    fill: color-mix(
      in srgb,
      var(--color-primary) 82%,
      var(--color-background)
    );
    stroke: color-mix(
      in srgb,
      var(--color-primary) 82%,
      var(--color-background)
    );
    stroke-width: 3;
    stroke-linecap: round;
    stroke-linejoin: round;
    opacity: 0.82;
  }

  .widget-footer {
    display: flex;
    align-items: center;
    justify-content: space-between;
    flex-wrap: wrap;
    gap: 12px;
    margin-top: 20px;
  }

  .widget-footer p {
    margin: 0;
    color: var(--color-text-muted);
    font-size: 12px;
  }

  .start-button {
    padding: 11px 18px;
    border: none;
    border-radius: 999px;
    background: var(--color-primary);
    color: var(--color-background);
    font-size: 13px;
    font-weight: 600;
  }

  .start-button:disabled {
    cursor: not-allowed;
    opacity: 0.65;
  }

  @media (max-width: 480px) {
    .scenic-card {
      min-height: 285px;
    }

    .timer-content {
      padding-top: 32px;
    }

    .timer-value {
      font-size: 52px;
    }
  }
</style>