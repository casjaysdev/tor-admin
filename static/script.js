// File: static/script.js
// Purpose: Combined logic for index.html and config.html functionality

let configData = [];

document.addEventListener('DOMContentLoaded', () => {
  if (document.getElementById('config-fields')) {
    renderConfigForm();
    hookConfigSave();
  }

  if (document.getElementById('amount')) {
    loadHiddenServices();
  }
});

// ========================
// CONFIG FORM (config.html)
// ========================
function renderConfigForm() {
  fetch('/api/options')
    .then((res) => res.json())
    .then((data) => {
      configData = data;
      const container = document.getElementById('config-fields');
      container.innerHTML = '';

      Object.keys(data).forEach((category) => {
        const section = document.createElement('div');
        section.innerHTML = `<h3 class="text-xl font-bold mb-2">${category}</h3>`;

        data[category].forEach((opt) => {
          const wrapper = document.createElement('div');
          wrapper.className = 'form-control w-full mb-2';

          const label = document.createElement('label');
          label.className = 'label justify-between';

          const labelText = document.createElement('span');
          labelText.className = 'label-text font-medium';
          labelText.textContent = `${opt.Name} (${opt.Type})`;

          const resetBtn = document.createElement('button');
          resetBtn.type = 'button';
          resetBtn.className = 'btn btn-xs btn-outline';
          resetBtn.textContent = 'Reset';
          resetBtn.onclick = () => {
            const field = document.getElementById(`opt-${opt.Name}`);
            field.value = opt.Default;
            if (field.type === 'checkbox') field.checked = opt.Default === '1';
          };

          label.appendChild(labelText);
          if (opt.Resettable) label.appendChild(resetBtn);

          const input = createInputField(opt);
          wrapper.appendChild(label);
          wrapper.appendChild(input);
          section.appendChild(wrapper);
        });

        container.appendChild(section);
      });
    });
}

function createInputField(opt) {
  const input = document.createElement('input');
  input.id = `opt-${opt.Name}`;
  input.name = opt.Name;
  input.placeholder = opt.Placeholder || '';
  input.className = 'input input-bordered w-full';

  switch (opt.InputType) {
    case 'number':
      input.type = 'number';
      input.value = opt.Default;
      break;
    case 'checkbox':
      input.type = 'checkbox';
      input.className = 'toggle';
      input.checked = opt.Default === '1';
      break;
    default:
      input.type = 'text';
      input.value = opt.Default;
  }

  return input;
}

function resetAllToDefaults() {
  configData &&
    Object.values(configData)
      .flat()
      .forEach((opt) => {
        const field = document.getElementById(`opt-${opt.Name}`);
        if (field) {
          if (field.type === 'checkbox') {
            field.checked = opt.Default === '1';
          } else {
            field.value = opt.Default;
          }
        }
      });
}

function hookConfigSave() {
  document.getElementById('config-form').addEventListener('submit', (e) => {
    e.preventDefault();
    const form = new FormData(e.target);
    const body = {};

    for (const [key, val] of form.entries()) {
      body[key] = val;
    }

    configData &&
      Object.values(configData)
        .flat()
        .forEach((opt) => {
          if (opt.InputType === 'checkbox') {
            const el = document.getElementById(`opt-${opt.Name}`);
            body[opt.Name] = el.checked ? '1' : '0';
          }
        });

    fetch('/api/torrc', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(body),
    })
      .then((res) => res.json())
      .then(() => alert('Configuration saved successfully.'))
      .catch(() => alert('Failed to save config'));
  });
}

// ========================
// BANDWIDTH + INDEX FEATURES
// ========================
function calculateBandwidth() {
  const amount = parseInt(document.getElementById('amount').value || '0');
  const unit = document.getElementById('unit').value;
  const interval = document.getElementById('interval').value;

  fetch('/api/bandwidth', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ amount, unit, interval }),
  })
    .then((res) => res.json())
    .then((data) => {
      document.getElementById('bw-result').innerText = `≈ ${data.human} (${data.bps} bytes/sec)`;
      renderGraph(amount, unit, interval);
    });
}

function renderGraph(amount, unit, interval) {
  const ctx = document.getElementById('bwChart')?.getContext('2d');
  if (!ctx) return;

  const seconds = {
    daily: 86400,
    weekly: 604800,
    monthly: 2592000,
  }[interval];

  const multiplier = {
    KB: 1024,
    MB: 1024 * 1024,
    GB: 1024 * 1024 * 1024,
    TB: 1024 * 1024 * 1024 * 1024,
  }[unit];

  const bps = (amount * multiplier) / seconds;
  const data = Array.from({ length: 24 }, () => Math.round(bps * 3600 * (Math.random() * 0.2 + 0.9)));

  if (window.bwChart) window.bwChart.destroy();
  window.bwChart = new Chart(ctx, {
    type: 'line',
    data: {
      labels: [...Array(24).keys()].map((h) => `${h}:00`),
      datasets: [
        {
          label: 'Projected Bandwidth (bytes/hour)',
          data,
          borderWidth: 2,
        },
      ],
    },
    options: {
      scales: { y: { beginAtZero: true } },
    },
  });
}

function loadHiddenServices() {
  fetch('/api/hidden')
    .then((res) => res.json())
    .then((data) => {
      const list = document.getElementById('onion-list');
      if (!list) return;
      list.innerHTML = '';
      data.onions.forEach((onion) => {
        const li = document.createElement('li');
        li.textContent = onion;
        list.appendChild(li);
      });
    });
}

function control(action) {
  fetch('/api/status?action=' + action)
    .then((res) => res.text())
    .then((data) => {
      const out = document.getElementById('svc-output');
      if (out) out.innerText = data;
    });
}
