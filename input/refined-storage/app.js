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
});
