import tippy from 'tippy.js';
import 'tippy.js/dist/tippy.css';
import 'tippy.js/themes/light.css';
import './style.scss';
import 'bootstrap';

window.addEventListener('load', () => {
    const tooltips = document.querySelectorAll('[data-tooltip-icon]');
    tooltips.forEach(tooltip => {
        tippy(tooltip, {
            content: '<img src="' + tooltip.getAttribute('data-tooltip-icon') + '" class="wiki-tooltip-icon">',
            theme: 'light',
            allowHTML: true,
            animation: null,
            placement: "right",
        })
    });

    tippy('[data-tippy-content]');

    if (localStorage.getItem('tab') === null) {
        localStorage.setItem('tab', 'Guides');
    }

    document.querySelector('#_sidebar_' + localStorage.getItem('tab')).style.display = 'block';
    document.querySelector('[data-tab="' + localStorage.getItem('tab') + '"]').classList.add('active');

    document.querySelectorAll('[data-toggle="tab"]').forEach(tab => {
        tab.addEventListener('click', () => {
            const name = tab.getAttribute('data-tab');

            document.querySelectorAll('.sidebar').forEach(item => item.style.display = 'none');
            document.querySelector('#_sidebar_' + name).style.display = 'block';

            localStorage.setItem('tab', name);
        });
    })
});
