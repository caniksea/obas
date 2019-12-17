$(document).ready(function () {

    /** Institution Location start here **/

    $("#province").change(function() {
        const provinceId = $(this).val();
        let districtElement = $('#district');
        let townElement = $('#town');
        getDropDownElement(townElement, 'Town');
        let districtDropDown = getDropDownElement(districtElement, "District");
        populateLocationDropDown(districtDropDown, provinceId);
    });

    $('#district').change(function() {
       // alert("in institution Id")
        const districtId = $(this).val();
        let townElement = $('#town');
        //let townDropDown = getDropDownElement(institutionElement, "Town");
        let townDropDown = getDropDownElement(townElement, "Town");
        populateLocationDropDown(townDropDown, districtId);
    });

    $("#institutionType").change(function() {

        const institutionTypeId = $(this).val();
        let institutionElement = $('#institution');
        let institutionDropDown = getDropDownElement(institutionElement, "Institution");
        populateInstitutionDropDownByType(institutionDropDown, institutionTypeId);
    });

    /** Institution Location ends here **/

});

function editForm(event) {
    var form = document.forms['siteEditForm'];
    form.elements["Id"].value = event.id;
    form.elements["Name"].value = event.name;
    form.elements["Description"].value = event.description;
}
function editLocationForm(app) {

    var form = document.forms['locationEditForm'];
    form.elements["LocationId"].value = app.locationId;
    form.elements["Name"].value = app.name;
    form.elements["Longitude"].value = app.longitude;
    form.elements["Latitude"].value = app.latitude;
}